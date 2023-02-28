---
title: "1-4 HAVING子句的力量"
date: 2023-02-27T21:44:54+08:00
draft: false
---

sql是面向集合的语言

# 寻找缺失的编号

SeqTbl

![](https://res.weread.qq.com/wrepub/epub_26211874_85)

用面向过程语言查询缺失的编号的过程

1. 对“连续编号”列按升序或者降序进行排序。
2. 循环比较每一行和下一行的编号。

将表整体看作一个集合

```sql
    -- 如果有查询结果，说明存在缺失的编号
    SELECT’存在缺失的编号’AS gap
      FROM SeqTbl
    HAVING COUNT(＊) <> MAX(seq);
```

执行结果

```sql
    gap
    ----------
    '存在缺失的编号’
```

上面的SQL语句里没有GROUP BY子句，此时整张表会被聚合为一行。这种情况下HAVING子句也是可以使用的。

再来查询一下缺失编号的最小值。

```sql
    -- 查询缺失编号的最小值
    SELECT MIN(seq + 1) AS gap
      FROM SeqTbl
     WHERE (seq+ 1) NOT IN ( SELECT seq FROM SeqTbl);
```

执行结果

```sql
    gap
    ---
      4
```

上面展示了通过SQL语句查询缺失编号的最基本的思路，然而这个查询还不够周全，并不能涵盖所有情况。例如，如果表SeqTbl里没有编号1，那么缺失编号的
最小值应该是1，但是这两条SQL语句都不能得出正确的结果

# 用HAVING子句进行子查询：求众数

Graduates（毕业生表）

![](https://res.weread.qq.com/wrepub/epub_26211874_91)

从这个例子可以看出，简单地求平均值有一个缺点，那就是很容易受到离群值（outlier）的影响。这种时候就必须使用更能准确反映出群体趋势的指标——
众数（mode）就是其中之一。

```sql
    --求众数的SQL语句(1)：使用谓词
    SELECT income, COUNT(＊) AS cnt
      FROM Graduates
     GROUP BY income
    HAVING COUNT(＊) >= ALL ( SELECT COUNT(＊)
                              FROM Graduates
                              GROUP BY income);
```

执行结果

```sql
    income  cnt
    ------  ---
    10000    3
    20000    3
```

[1-3节]({{< ref "posts/books/SQL进阶教程/第1章-神奇的SQL/1-3-三值逻辑和NULL.md#qualified-predicates-and-null" >}})提到过ALL谓
词用于NULL或空集时会出现问题，可以用极值函数来代替。

```sql
    --求众数的SQL语句(2)：使用极值函数
    SELECT income, COUNT(＊) AS cnt
      FROM Graduates
     GROUP BY income
    HAVING COUNT(＊) >=  ( SELECT MAX(cnt)
                            FROM ( SELECT COUNT(＊) AS cnt
                                    FROM Graduates
                                  GROUP BY income) TMP ) ;
```

# 用HAVING子句进行自连接：求中位数

用面向集合的方式，来查询位于集合正中间的元素。

- 将集合里的元素按照大小分为上半部分和下半部分两个子集，同时让这2个子集共同拥有集合正中间的元素。

中位数求法的思路

![](https://res.weread.qq.com/wrepub/epub_26211874_94)

```sql
    --求中位数的SQL语句：在HAVING子句中使用非等值自连接
    SELECT AVG(DISTINCT income)
      FROM (SELECT T1.income
              FROM Graduates T1, Graduates T2
            GROUP BY T1.income

            --S1的条件
      HAVING SUM(CASE WHEN T2.income >= T1.income THEN 1 ELSE 0 END)
                >= COUNT(＊) / 2
            --S2的条件
        AND SUM(CASE WHEN T2.income <= T1.income THEN 1 ELSE 0 END)
                >= COUNT(＊) / 2 ) TMP;
```

# 查询不包含NULL的集合

count(*)和count(column)的区别

- 第一个是性能上的区别；第二个是COUNT（＊）可以用于NULL，而COUNT（列名）与其他聚合函数一样，要先排除掉NULL的行再进行统计。
- 第二个区别也可以这么理解：COUNT（＊）查询的是所有行的数目，而COUNT（列名）查询的则不一定是。

一张全是NULL的表NullTbl

![](https://res.weread.qq.com/wrepub/epub_26211874_96)

```sql
    --在对包含NULL的列使用时，COUNT（＊）和COUNT（列名）的查询结果是不同的
    SELECT COUNT(＊), COUNT(col_1)
      FROM NullTbl;
```

执行结果

```sql
    count(＊)   count(col_1)
    --------   ------------
          3               0
```

Students

![](https://res.weread.qq.com/wrepub/epub_26211874_97)

所有学生都提交了报告的学院有哪些

![](https://res.weread.qq.com/wrepub/epub_26211874_98)

```sql
    --查询“提交日期”列内不包含NULL的学院(1)：使用COUNT函数
    SELECT dpt
      FROM Students
     GROUP BY dpt
    HAVING COUNT(＊) = COUNT(sbmt_date);
```

执行结果

```sql
    dpt
    --------
    理学院
    经济学院
```

使用CASE表达式也可以实现同样的功能

```sql
    --查询“提交日期”列内不包含NULL的学院(2)：使用CASE表达式
    SELECT dpt
      FROM Students
     GROUP BY dpt
    HAVING COUNT(＊) = SUM(CASE WHEN sbmt_date IS NOT NULL
                              THEN 1
                              ELSE 0 END);
```

# 用关系除法运算进行购物篮分析

我们假设有这样两张表：全国连锁折扣店的商品表Items，以及各个店铺的库存管理表ShopItems。

Items

![](https://res.weread.qq.com/wrepub/epub_26211874_100)

ShopItems

![](https://res.weread.qq.com/wrepub/epub_26211874_101)

查询囊括了表Items中所有商品的店铺

```sql
    --查询啤酒、纸尿裤和自行车同时在库的店铺：错误的SQL语句
    SELECT DISTINCT shop
      FROM ShopItems
     WHERE item IN (SELECT item FROM Items);
```

执行结果

```sql
    shop
    ----
    仙台
    东京
    大阪
```

```sql
    --查询啤酒、纸尿裤和自行车同时在库的店铺：正确的SQL语句
    SELECT SI.shop
      FROM ShopItems SI, Items I
     WHERE SI.item = I.item
     GROUP BY SI.shop
    HAVING COUNT(SI.item) = (SELECT COUNT(item) FROM Items);
```

执行结果

```sql
    shop
    ----
    仙台
    东京
```

**请注意，如果把HAVING子句改成HAVING COUNT(SI.item)=COUNT(I.item)，结果就不对了。**

```sql
    -- COUNT(I.item)的值已经不一定是3了
    SELECT SI.shop, COUNT(SI.item), COUNT(I.item)
      FROM ShopItems SI, Items I
     WHERE SI.item = I.item
     GROUP BY SI.shop;
```

执行结果

```sql
    shop   COUNT(SI.item)   COUNT(I.item)
    -----  ---------------  --------------
    仙台                   3                3
    东京                   3                3
    大阪                   2                2
```

如何排除掉仙台店（仙台店的仓库中存在“窗帘”，但商品表里没有“窗帘”），让结果里只出现东京店。

- 这类问题被称为“精确关系除法”（exact relational division），即只选择没有剩余商品的店铺
- 与此相对，前一个问题被称为“带余除法”（division with a remainder）。

```sql
    --精确关系除法运算：使用外连接和COUNT函数
    SELECT SI.shop
      FROM ShopItems SI LEFT OUTER JOIN Items I
        ON SI.item=I.item
     GROUP BY SI.shop
    HAVING COUNT(SI.item) = (SELECT COUNT(item) FROM Items)    --条件1
      AND COUNT(I.item)  = (SELECT COUNT(item) FROM Items);   --条件2
```

执行结果

```sql
    shop
    ----
     东京
```

表ShopItems和表Items外连接后的结果

![](https://res.weread.qq.com/wrepub/epub_26211874_103)

一般来说，使用外连接时，大多会用商品表Items作为主表进行外连接操作，而这里颠倒了一下主从关系，表使用ShopItems作为了主表，这一点比较有趣。

# 本节小结

1. 表不是文件，记录也没有顺序，所以SQL不进行排序。
2. SQL不是面向过程语言，没有循环、条件分支、赋值操作。
3. SQL通过不断生成子集来求得目标集合。SQL不像面向过程语言那样通过画流程图来思考问题，而是通过画集合的关系图来思考。
4. GROUP BY子句可以用来生成子集。
5. WHERE子句用来调查集合元素的性质，而HAVING子句用来调查集合本身的性质。
