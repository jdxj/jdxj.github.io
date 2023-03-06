---
title: "1-7 用SQL进行集合运算"
date: 2023-03-02T21:38:49+08:00
draft: false
---

# 导入篇：集合运算的几个注意事项

注意事项1: SQL能操作具有重复行的集合，可以通过可选项ALL来支持

- 一般的集合论不允许集合中存在重复元素
- 关系数据库允许存在重复行, 称为多重集合（multiset, bag）
- SQL的集合运算符也提供了允许重复和不允许重复的两种用法。如果直接使用UNION或INTERSECT，结果里就不会出现重复的行。如果想在结果里留下重复行，
  可以加上可选项ALL，写作UNION ALL。
- **集合运算符为了排除掉重复行，默认地会发生排序，而加上可选项ALL之后，就不会再排序，所以性能会有提升。**

注意事项2：集合运算符有优先级

- 标准SQL规定，INTERSECT比UNION和EXCEPT优先级更高。

注意事项3：各个DBMS提供商在集合运算的实现程度上参差不齐

- SQL Server从2005版开始支持INTERSECT和EXCEPT，而MySQL还都不支持（包含在“中长期计划”里）。
- 还有像Oracle这样，实现了EXCEPT功能但却命名为MINUS的数据库。

注意事项4：除法运算没有标准定义

四则运算里的和（UNION）、差（EXCEPT）、积（CROSS JOIN）都被引入了标准SQL。但是很遗憾，商（DIVIDE BY）因为各种原因迟迟没能标准化

# 比较表和表：检查集合相等性之基础篇

“相等”指的是行数和列数以及内容都相同

名字不同但内容相同的两张表

tbl_A

![](https://res.weread.qq.com/wrepub/epub_26211874_176)

tbl_B

![](https://res.weread.qq.com/wrepub/epub_26211874_177)

如果这个查询的结果与tbl_A及tbl_B的行数一致，则两张表是相等的

- **需要事先确认表tbl_A和表tbl_B的行数是一样的**

```sql
    SELECT COUNT(＊) AS row_cnt
      FROM ( SELECT ＊
              FROM tbl_A
            UNION
            SELECT ＊
              FROM tbl_B ) TMP;
```

执行结果

```sql
    row_cnt
    -------
          3
```

如果集合运算符里不加上可选项ALL，那么重复行就会被排除掉。因此，如果表tbl_A和表tbl_B是相等的，**排除掉重复行后，两个集合是完全重合的**。

![](https://res.weread.qq.com/wrepub/epub_26211874_178)

key列为B的一行数据不同：结果会变为4

tbl_A
![](https://res.weread.qq.com/wrepub/epub_26211874_179)

tbl_B

![](https://res.weread.qq.com/wrepub/epub_26211874_180)

前面的SQL语句可以用于包含NULL数据的表，而且不需要指定列数、列名和数据类型等就能使用

对于任意的表S，都有下面的公式成立。

```sql
S UNION S = S
```

同一个集合无论加多少次结果都相同

```sql
S UNION S UNION S UNION S …… UNION S = S
```

**UNION的这个优雅而强大的幂等性只适用于数学意义上的集合，对SQL中有重复数据的多重集合是不适用的。**

# 比较表和表：检查集合相等性之进阶篇

在集合论里，判定两个集合是否相等时

- (A ∩ B )且(A ∩ B) ⇔ (A = B) 如果集合A包含集合B，且集合B包含集合A，则集合A和集合B相等
- (A ∪ B ) = (A ∩ B) ⇔ (A = B)

除了UNION之外，另一个具有幂等性的运算符就是INTERSECT。

![](https://res.weread.qq.com/wrepub/epub_26211874_183)

如果A = B，则(A UNION B) EXCEPT (A INTERSECT B)的结果是空集

```sql
    --两张表相等时返回“相等”，否则返回“不相等”
    SELECT CASE WHEN COUNT(＊) = 0
                THEN ’相等’
                ELSE’不相等’END AS result
      FROM ((SELECT ＊ FROM  tbl_A
            UNION
            SELECT ＊ FROM  tbl_B)
            EXCEPT
            (SELECT ＊ FROM  tbl_A
            INTERSECT
            SELECT ＊ FROM  tbl_B)) TMP;
```

这条SQL语句与上一部分中的SQL语句具有同样的优点，也不需要知道列名和列数，还可以用于包含NULL的表，而且，这个改进版连事先查询两张表的行数这种
准备工作也不需要了。

虽然功能改进了，却也带来了一些缺陷。由于这里需要进行4次排序（3次集合运算加上1次DISTINCT），所以性能会有所下降（不过这条SQL语句也不需要频繁
执行，所以这点缺陷也不是不能容忍）。

**因为这里使用了INTERSECT和EXCEPT，所以目前这条SQL语句不能在MySQL里执行。**

把不同的行输出

```sql
    --用于比较表与表的diff
    (SELECT ＊ FROM  tbl_A
     EXCEPT
     SELECT ＊ FROM  tbl_B)
    UNION ALL
    (SELECT ＊ FROM  tbl_B
     EXCEPT
     SELECT ＊ FROM  tbl_A);
```

执行结果

```sql
    key   col_1   col_2   col_3
    ---   -----   -----   -----
    B          0        7        9
    B          0        7        8
```

- 因为A－B和B－A之间不可能有交集，所以合并这两个结果时使用UNION ALL也没有关系。
- 在A和B一方包含另一方时，这条SQL语句也是成立的（这时A－B或者B－A有一个会是空集）。
- 需要注意的是，在SQL中，括号决定了运算的先后顺序，非常重要，如果去掉括号，结果就会不正确。

# 用差集实现关系除法运算

实现除法的方法

- 嵌套使用NOT EXISTS。
- 使用HAVING子句转换成一对一关系。
- 把除法变成减法。

本节将介绍一下第三种方法。

两张员工技术信息管理表

Skills

![](https://res.weread.qq.com/wrepub/epub_26211874_185)

EmpSkills

![](https://res.weread.qq.com/wrepub/epub_26211874_186)

问题是，从表EmpSkills中找出精通表Skills中所有技术的员工。也就是说，答案是相田和神崎。

```sql
    --用求差集的方法进行关系除法运算（有余数）
    SELECT DISTINCT emp
      FROM EmpSkills ES1
     WHERE NOT EXISTS
            (SELECT skill
              FROM Skills
            EXCEPT
            SELECT skill
              FROM EmpSkills ES2
              WHERE ES1.emp = ES2.emp);
```

执行结果

```sql
    emp
    ---
    相田
    神崎
```

从需求的技术的集合中减去每个员工自己的技术的集合，如果结果是空集，则说明该员工具备所有的需求的技术，否则说明该员工不具备某些需求的技术。

相田

![](https://res.weread.qq.com/wrepub/epub_26211874_187)

平井

![](https://res.weread.qq.com/wrepub/epub_26211874_188)

# 寻找相等的子集

供应商及其经营的零件的表

SupParts

![](https://res.weread.qq.com/wrepub/epub_26211874_193)

我们需要求的是，经营的零件在种类数和种类上都完全相同的供应商组合。

SQL并没有提供任何用于检查集合的包含关系或者相等性的谓词。IN谓词只能用来检查元素是否属于某个集合（∈），而不能检查集合是否是某个集合的子集
（∪）。

首先，我们来生成供应商的全部组合。

```sql
    --生成供应商的全部组合
    SELECT SP1.sup AS s1, SP2.sup AS s2
      FROM SupParts SP1, SupParts SP2
     WHERE SP1.sup < SP2.sup
     GROUP BY SP1.sup, SP2.sup;
```

执行结果

```sql
    s1    s2
    ----  ----
    A     B
    A     C
    A     D
      ：
      ：
      ：
    D     E
    E     F
```

接下来，我们检查一下这些供应组合是否满足以下公式：[插图]且[插图]。这个公式等价于下面两个条件。

- 条件1：两个供应商都经营同种类型的零件
- 条件2：两个供应商经营的零件种类数相同（即存在一一映射）

条件1只需要简单地按照“零件”列进行连接，而条件2需要用COUNT函数来描述。

```sql
    SELECT SP1.sup AS s1, SP2.sup AS s2
      FROM SupParts SP1, SupParts SP2
     WHERE SP1.sup < SP2.sup                  --生成供应商的全部组合
      AND SP1.part = SP2.part                --条件1：经营同种类型的零件
     GROUP BY SP1.sup, SP2.sup
    HAVING COUNT(＊) = (SELECT COUNT(＊)       --条件2：经营的零件种类数相同
                        FROM SupParts SP3
                        WHERE SP3.sup = SP1.sup)
      AND COUNT(＊) = (SELECT COUNT(＊)
                        FROM SupParts SP4
                        WHERE SP4.sup = SP2.sup);
```

> 因为要比较任意两个供应商的零件, 所以用笛卡儿积

执行结果

```sql
    s1    s2
    ----  ----
    A     C
    B     D
```

# 用于删除重复行的高效SQL

![](https://res.weread.qq.com/wrepub/epub_26211874_198)

1-2节介绍的解法是使用关联子查询

```sql
    --删除重复行：使用关联子查询
    DELETE FROM Products
     WHERE rowid < ( SELECT MAX(P2.rowid)
                      FROM Products P2
                      WHERE Products.name  = P2. name
                        AND Products.price = P2.price ) ;
```

上面这条语句的思路是，按照“商品名，价格”的组合汇总后，求出每个组合的最大rowid，然后把其余的行都删除掉。

假设表中加上了“rowid”列

![](https://res.weread.qq.com/wrepub/epub_26211874_199)

```sql
    --用于删除重复行的高效SQL语句(1)：通过EXCEPT求补集
    DELETE FROM Products
     WHERE rowid IN ( SELECT rowid           --全部rowid
                        FROM Products
                      EXCEPT                 --减去
                      SELECT MAX(rowid)     --要留下的rowid
                        FROM Products
                      GROUP BY name, price) ;
```

使用EXCEPT求补集的逻辑如下面的图表所示。

![](https://res.weread.qq.com/wrepub/epub_26211874_200)

把EXCEPT改写成NOT IN也是可以实现的。

```sql
    --删除重复行的高效SQL语句(2)：通过NOT IN求补集
    DELETE FROM Products
     WHERE rowid NOT IN ( SELECT MAX(rowid)
                            FROM Products
                          GROUP BY name, price);
```

