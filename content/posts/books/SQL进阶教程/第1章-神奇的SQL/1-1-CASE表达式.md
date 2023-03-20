---
title: "1-1 CASE表达式"
date: 2023-02-25T16:27:54+08:00
---

# CASE表达式概述

CASE表达式的写法

```sql
    --简单CASE表达式
    CASE sex
      WHEN '1' THEN ’男’
      WHEN '2' THEN ’女’
    ELSE ’其他’ END

    --搜索CASE表达式
    CASE WHEN sex ='1'THEN’男’
        WHEN sex ='2'THEN’女’
    ELSE ’其他’ END
```

剩余的WHEN子句被忽略的写法示例

```sql
    --例如，这样写的话，结果里不会出现“第二”
    CASE WHEN col_1 IN ('a', 'b') THEN’第一’
        WHEN col_1 IN ('a')     THEN’第二’
    ELSE ’其他’ END
```

注意

- 注意事项1：各分支返回的数据类型要一致, 否则报错
- 注意事项2：不要忘了写END
- 注意事项3：养成写ELSE子句的习惯
  - 不写可能会造成“语法没有错误，结果却不对”这种不易追查原因的麻烦
  - 养成这样的习惯后，我们从代码上就可以清楚地看到这种条件下会生成NULL，而且将来代码有修改时也能减少失误。

# 将已有编号方式转换为新的方式并统计

统计数据源表PopTbl

![](https://res.weread.qq.com/wrepub/epub_26211874_9)

统计结果

![](https://res.weread.qq.com/wrepub/epub_26211874_11)

用县名（pref_name）代替编号作为GROUP BY的列

```sql
    --把县编号转换成地区编号(1)
    SELECT  CASE pref_name
                    WHEN ’德岛’ THEN ’四国’
                    WHEN ’香川’ THEN ’四国’
                    WHEN ’爱媛’ THEN ’四国’
                    WHEN ’高知’ THEN ’四国’

                    WHEN ’福冈’ THEN ’九州’
                    WHEN ’佐贺’ THEN ’九州’
                    WHEN ’长崎’ THEN ’九州’
            ELSE’其他’END AS district,
            SUM(population)
      FROM  PopTbl
      GROUP BY CASE pref_name
                    WHEN ’德岛’ THEN ’四国’
                    WHEN ’香川’ THEN ’四国’
                    WHEN ’爱媛’ THEN ’四国’
                    WHEN ’高知’ THEN ’四国’
                    WHEN ’福冈’ THEN ’九州’
                    WHEN ’佐贺’ THEN ’九州’
                    WHEN ’长崎’ THEN ’九州’
              ELSE ’其他’ END;
```

将数值按照适当的级别进行分类统计

```sql
    --按人口数量等级划分都道府县
    SELECT  CASE WHEN population <  100 THEN'01'
                WHEN population >= 100 AND population < 200  THEN'02'
                WHEN population >= 200 AND population < 300  THEN'03'
                WHEN population >= 300 THEN'04'
            ELSE NULL END AS pop_class,
            COUNT(＊) AS cnt
      FROM  PopTbl
     GROUP BY CASE WHEN population <  100 THEN'01'
                  WHEN population >= 100 AND population < 200  THEN'02'
                  WHEN population >= 200 AND population < 300  THEN'03'
                  WHEN population >= 300 THEN'04'
              ELSE NULL END;

    pop_class  cnt
    --------- ----
    01            1
    02            3
    03            3
    04            2
```

上面两种方式都要在select和group by中写相同的case, 在变更时可能会忘记同步, 下面的写法更方便

![](https://res.weread.qq.com/wrepub/epub_26211874_12)

- 严格来说，这种写法是违反标准SQL的规则的。因为GROUP BY子句比SELECT语句先执行，所以在GROUP BY子句中引用在SELECT子句里定义的别称是不被允
  许的。事实上，在Oracle、DB2、SQL Server等数据库里采用这种写法时就会出错。
- 不过也有支持这种SQL语句的数据库，例如在PostgreSQL和MySQL中，这个查询语句就可以顺利执行。这是因为，这些数据库在执行查询语句时，会先对
  SELECT子句里的列表进行扫描，并对列进行计算。

# 用一条SQL语句进行不同条件的统计

统计源表PopTbl2

![](https://res.weread.qq.com/wrepub/epub_26211874_14)

统计结果

![](https://res.weread.qq.com/wrepub/epub_26211874_15)

通常的做法是写两个sql, 可能然后再用union合并

```sql
    -- 男性人口
    SELECT pref_name,
          SUM(population)
      FROM PopTbl2
     WHERE sex ='1'
     GROUP BY pref_name;


    -- 女性人口
    SELECT pref_name,
          SUM(population)
      FROM PopTbl2
     WHERE sex ='2'
     GROUP BY pref_name;
```

用一条sql实现

```sql
    SELECT pref_name,
          --男性人口
          SUM( CASE WHEN sex ='1'THEN population ELSE 0 END) AS cnt_m,
          --女性人口
          SUM( CASE WHEN sex ='2'THEN population ELSE 0 END) AS cnt_f
      FROM  PopTbl2
     GROUP BY pref_name;
```

# 用CHECK约束定义多个列的条件关系

假设某公司规定“女性员工的工资必须在20万日元以下”，而在这个公司的人事表中，这条无理的规定是使用CHECK约束来描述的

```sql
    CONSTRAINT check_salary CHECK
              ( CASE
                    WHEN sex ='2' THEN 
                        CASE
                            WHEN salary <= 200000 THEN 1
                            ELSE 0
                        END
                    ELSE 1
                END = 1 )
```

用逻辑与改写的CHECK约束如下所示。

```sql
    CONSTRAINT check_salary CHECK
              ( sex ='2'AND salary <= 200000 )
```

这两个约束的程序行为不一样(第二个如果是男雇员也返回false了, 不符合命题)

逻辑与和蕴含式的真值表

- U表示unknown

![](https://res.weread.qq.com/wrepub/epub_26211874_17)

# 在UPDATE语句里进行条件分支

Salaries

![](https://res.weread.qq.com/wrepub/epub_26211874_19)

假设现在需要根据以下条件对该表的数据进行更新。

1. 对当前工资为30万日元以上的员工，降薪10%。
2. 对当前工资为25万日元以上且不满28万日元的员工，加薪20%。

按照这些要求更新完的数据应该如下表所示。

![](https://res.weread.qq.com/wrepub/epub_26211874_20)

乍一看，分别执行下面两个UPDATE操作好像就可以做到，但这样的结果却是不正确的。

```sql
    --条件1
    UPDATE Salaries
      SET salary = salary ＊ 0.9
     WHERE salary >= 300000;

    --条件2
    UPDATE Salaries
      SET salary = salary ＊ 1.2
     WHERE salary >= 250000 AND salary < 280000;
```

正确的写法

- 注意最后的else必须写, 否则返回null

```sql
    --用CASE表达式写正确的更新操作
    UPDATE Salaries
      SET salary = CASE WHEN salary >= 300000
                        THEN salary ＊ 0.9
                        WHEN salary >= 250000 AND salary < 280000
                        THEN salary ＊ 1.2
                        ELSE salary END;
```

调换主键值的方便写法

SomeTable

![](https://res.weread.qq.com/wrepub/epub_26211874_22)

不使用case

```sql
    --1．将a转换为中间值d
    UPDATE SomeTable
      SET p_key ='d'
     WHERE p_key ='a';


    --2．将b调换为a
    UPDATE SomeTable
      SET p_key ='a'

      WHERE p_key ='b';


     --3．将d调换为b
     UPDATE SomeTable
        SET p_key ='b'
      WHERE p_key ='d';
```

使用case

- 适用于Oracle, DB2, SQL Server
- 不适用于PostgreSQL, MySQWL

```sql
    --用CASE表达式调换主键值
    UPDATE SomeTable
      SET p_key = CASE WHEN p_key ='a'
                        THEN 'b'
                        WHEN p_key ='b'
                        THEN 'a'
                        ELSE p_key END
     WHERE p_key IN ('a', 'b');
```

# 表之间的数据匹配

课程一览 CourseMaster

![](https://res.weread.qq.com/wrepub/epub_26211874_25)

开设的课程 OpenCourses

![](https://res.weread.qq.com/wrepub/epub_26211874_26)

我们要用这两张表来生成下面这样的交叉表

```
    course_name   6月   7月   8月
    -----------  ----  ----  ----
    会计入门         ○    ×     ×
    财务知识         ×    ×    ○
    簿记考试         ○    ×     ×
    税务师           ○    ○    ○
```

```sql
    --表的匹配：使用IN谓词
    SELECT course_name,
          CASE WHEN course_id IN
                        (SELECT course_id FROM OpenCourses
                          WHERE month = 200706) THEN'○'
                ELSE'×'END AS "6月",
          CASE WHEN course_id IN
                        (SELECT course_id FROM OpenCourses
                          WHERE month = 200707) THEN'○'
                ELSE'×'END AS "7月",
          CASE WHEN course_id IN
                        (SELECT course_id FROM OpenCourses
                          WHERE month = 200708) THEN'○'
                ELSE'×'END  AS "8月"
      FROM CourseMaster;


    --表的匹配：使用EXISTS谓词
    SELECT CM.course_name,
          CASE WHEN EXISTS
                        (SELECT course_id FROM OpenCourses OC
                          WHERE month = 200706

                              AND OC.course_id = CM.course_id) THEN'○'
                  ELSE'×'END AS "6月",
              CASE WHEN EXISTS
                          (SELECT course_id FROM OpenCourses OC
                            WHERE month = 200707
                              AND OC.course_id = CM.course_id) THEN'○'
                  ELSE'×'END AS "7月",
              CASE WHEN EXISTS
                          (SELECT course_id FROM OpenCourses OC
                            WHERE month = 200708
                              AND OC.course_id = CM.course_id) THEN'○'
                  ELSE'×'END  AS "8月"
        FROM CourseMaster CM;
```

无论使用IN还是EXISTS，得到的结果是一样的，但从性能方面来说，EXISTS更好。通过EXISTS进行的子查询能够用到“month, course_id”这样的主键索引，
因此尤其是当表OpenCourses里数据比较多的时候更有优势。

# 在CASE表达式中使用聚合函数

假设这里有一张显示了学生及其加入的社团的一览表。如表StudentClub所示，这张表的主键是“学号、社团ID”，存储了学生和社团之间多对多的关系。

StudentClub

![](https://res.weread.qq.com/wrepub/epub_26211874_28)

我们按照下面的条件查询这张表里的数据。

1. 获取只加入了一个社团的学生的社团ID。
2. 获取加入了多个社团的学生的主社团ID。

条件1的SQL

```sql
    --条件1：选择只加入了一个社团的学生
    SELECT std_id, MAX(club_id) AS main_club
      FROM StudentClub
     GROUP BY std_id
    HAVING COUNT(＊) = 1;
```

执行结果1

```sql
    std_id   main_club
    ------   ----------
    300       4
    400       5
    500       6
```

条件2的SQL

```sql
    --条件2：选择加入了多个社团的学生
    SELECT std_id, club_id AS main_club
      FROM StudentClub
     WHERE main_club_flg ='Y';
```

执行结果2

```sql
    std_id  main_club
    ------  ----------
    100     1
    200     3
```

如果使用CASE表达式，下面这一条SQL语句就可以了

```sql
    SELECT  std_id,
            CASE WHEN COUNT(＊) = 1  --只加入了一个社团的学生
                THEN MAX(club_id)
                ELSE MAX(CASE WHEN main_club_flg ='Y'
                              THEN club_id
                              ELSE NULL END)
            END AS main_club
      FROM StudentClub
     GROUP BY std_id;

    std_id   main_club
    ------   ----------
    100       1
    200       3
    300       4
    400       5
    500       6
```

# 本节小结

最后说一点细节的东西。CASE表达式经常会因为同VB和C语言里的CASE“语句”混淆而被叫作CASE语句。但是准确来说，它并不是语句，而是和1+1或者a/b一样
属于表达式的范畴。结束符END确实看起来像是在标记一连串处理过程的终结，所以初次接触CASE表达式的人容易对这一点感到困惑。“表达式”和“语句”的名称
区别恰恰反映了两者在功能处理方面的差异。

作为表达式，CASE表达式在执行时会被判定为一个固定值，因此它可以写在聚合函数内部；也正因为它是表达式，所以还可以写在SELECE子句、GROUP BY
子句、WHERE子句、ORDER BY子句里。简单点说，在能写列名和常量的地方，通常都可以写CASE表达式。从这个意义上来说，与CASE表达式最接近的不是面向
过程语言里的CASE语句，而是Lisp和Scheme等函数式语言里的case和cond这样的条件表达式。
