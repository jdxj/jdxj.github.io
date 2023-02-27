---
title: "1-3 三值逻辑和NULL"
date: 2023-02-26T12:06:38+08:00
draft: false
---

SQL是三值逻辑(three-valued logic)

# 理论篇

两种NULL、三值逻辑还是四值逻辑

- 未知(unknown). 以“不知道戴墨镜的人眼睛是什么颜色”这种情况为例，这个人的眼睛肯定是有颜色的，但是如果他不摘掉眼镜，别人就不知道他的眼睛是什
  么颜色。
- 不适用(not applicable, inapplicable). 而“不知道冰箱的眼睛是什么颜色”则属于“不适用”。因为冰箱根本就没有眼睛

关系数据库中“丢失的信息”的分类

![](https://res.weread.qq.com/wrepub/epub_26211874_66)

四值逻辑真值表

![](https://res.weread.qq.com/wrepub/epub_26211874_68)

为什么必须写成“IS NULL”，而不是“＝NULL”

- 对NULL使用比较谓词后得到的结果总是unknown。

```sql
    --以下的式子都会被判为 unknown
    1 = NULL
    2 > NULL
    3 < NULL
    4 <> NULL
    NULL = NULL
```

为什么对NULL使用比较谓词后得到的结果永远不可能为真呢？

- 因为NULL既不是值也不是变量。NULL只是一个表示“没有值”的标记，而比较谓词只适用于值。因此，对并非值的NULL使用比较谓词本来就是没有意义的

unknown、第三个真值

- 书中用unknown表示三值逻辑的第三个值
- 书中用UNKNOWN表示NULL

```sql
    --这个是明确的真值的比较
    unknown = unknown → true

    --这个相当于NULL = NULL
    UNKNOWN = UNKNOWN → unknown
```

三值逻辑的真值表（NOT）

![](https://res.weread.qq.com/wrepub/epub_26211874_70)

三值逻辑的真值表（AND）

![](https://res.weread.qq.com/wrepub/epub_26211874_71)

三值逻辑的真值表（OR）

![](https://res.weread.qq.com/wrepub/epub_26211874_72)

记忆方法

- AND的情况： false ＞ unknown ＞ true
- OR的情况： true ＞ unknown ＞ false
- 优先级高的真值会决定计算结果。例如true AND unknown，因为unknown的优先级更高，所以结果是unknown。而true OR unknown的话，因为true优
  先级更高，所以结果是true。
- 特别需要记住的是，当AND运算中包含unknown时，结果肯定不会是true（反之，如果AND运算结果为true，则参与运算的双方必须都为true）。

# 实践篇

## 1. 比较谓词和NULL(1)：排中律不成立

“把命题和它的否命题通过‘或者’连接而成的命题全都是真命题”这个命题在二值逻辑中被称为排中律（Law of Excluded Middle）。顾名思义，排中律就是
指不认可中间状态，对命题真伪的判定黑白分明，是古典逻辑学的重要原理。“是否承认这一原理”被认为是古典逻辑学和非古典逻辑学的分界线。

如果排中律在SQL里也成立，那么下面的查询应该能选中表里的所有行。

```sql
    --查询年龄是20岁或者不是20岁的学生
    SELECT ＊
      FROM Students
     WHERE age = 20
        OR age <> 20;
```

遗憾的是，在SQL的世界里，排中律是不成立的。假设表Students里的数据如下所示。

Students

![](https://res.weread.qq.com/wrepub/epub_26211874_74)

这条SQL语句无法查询到约翰，因为约翰年龄不详。

SQL语句的查询结果里只有判断结果为true的行。要想让约翰出现在结果里，需要添加下面这样的“第3个条件”。

```sql
    --添加第3个条件：年龄是20岁，或者不是20岁，或者年龄未知
    SELECT ＊
      FROM Students
     WHERE age = 20
        OR age <> 20
        OR age IS NULL;
```

## 2. 比较谓词和NULL(2):CASE表达式和NULL

```sql
    --col_1为1时返回○、为NULL时返回×的CASE表达式？
    CASE col_1
      WHEN 1     THEN'○'
      WHEN NULL  THEN'×'
    END
```

这个CASE表达式一定不会返回×。这是因为，第二个WHEN子句是col_1 = NULL的缩写形式。这个式子的真值永远是unknown。

正确写法

```sql
    CASE WHEN col_1 = 1 THEN'○'
        WHEN col_1 IS NULL THEN'×'
     END
```

## 3. NOT IN和NOT EXISTS不是等价的

Class_A

![](https://res.weread.qq.com/wrepub/epub_26211874_76)

Class_B

![](https://res.weread.qq.com/wrepub/epub_26211874_77)

考虑一下如何根据这两张表查询“与B班住在东京的学生年龄不同的A班学生”。

- 希望查询到的是拉里和伯杰。因为布朗与齐藤年龄相同，所以不是我们想要的结果。
- 如果单纯地按照这个条件去实现，则SQL语句如下所示。

```sql
    --查询与B班住在东京的学生年龄不同的A班学生的SQL语句？
    SELECT ＊
      FROM Class_A
     WHERE age NOT IN ( SELECT age
                          FROM Class_B
                        WHERE city =’东京’);
```

上面的sql查不到任何数据

分析

```sql
--1．执行子查询，获取年龄列表
SELECT ＊
  FROM Class_A
 WHERE age NOT IN (22, 23, NULL);

--2．用NOT和IN等价改写NOT IN
SELECT ＊
FROM Class_A
WHERE NOT age IN (22, 23, NULL);

--3．用OR等价改写谓词IN
SELECT ＊
FROM Class_A
WHERE NOT ( (age = 22) OR (age = 23) OR (age = NULL) );

--4．使用德·摩根定律等价改写
SELECT ＊
FROM Class_A
WHERE NOT (age = 22) AND NOT(age = 23) AND NOT (age = NULL);

--5．用<>等价改写NOT和=
SELECT ＊
FROM Class_A
WHERE (age <> 22) AND (age <> 23) AND (age <> NULL);

--6．对NULL使用<>后，结果为unknown
SELECT ＊
FROM Class_A
WHERE (age <> 22) AND (age <> 23) AND unknown;

--7．如果AND运算里包含unknown，则结果不为true（参考“理论篇”中的矩阵）
SELECT ＊
FROM Class_A
WHERE false或unknown;
```

**如果NOT IN子查询中用到的表里被选择的列中存在NULL，则SQL语句整体的查询结果永远是空。**

为了得到正确的结果，我们需要使用EXISTS谓词。

```sql
    --正确的SQL语句：拉里和伯杰将被查询到
    SELECT ＊
      FROM Class_A  A
     WHERE NOT EXISTS ( SELECT ＊
                          FROM Class_B B
                        WHERE A.age = B.age
                          AND B.city = ’东京’);
```

执行结果

```sql
    name   age   city
    -----  ----  ----
    拉里    19    埼玉
    伯杰    21    千叶
```

分析

```sql
--1．在子查询里和NULL进行比较运算
SELECT ＊
  FROM Class_A A
 WHERE NOT EXISTS ( SELECT ＊
                      FROM Class_B B
                    WHERE A.age = NULL
                      AND B.city =’东京’);

--2．对NULL使用“=”后，结果为 unknown
SELECT ＊
FROM Class_A A
WHERE NOT EXISTS ( SELECT ＊
                   FROM Class_B B
                   WHERE unknown
                     AND B.city =’东京’);

--3．如果AND运算里包含unknown，结果不会是true
SELECT ＊
FROM Class_A A
WHERE NOT EXISTS ( SELECT ＊
                   FROM Class_B B
                   WHERE false或unknown);

--4．子查询没有返回结果，因此相反地，NOT EXISTS为true
SELECT ＊
FROM Class_A A
WHERE true;
```

## 4. 限定谓词和NULL

any与in等价

Class_A

![](https://res.weread.qq.com/wrepub/epub_26211874_78)

Class_B

![](https://res.weread.qq.com/wrepub/epub_26211874_79)

思考一下用于查询“比B班住在东京的所有学生年龄都小的A班学生”的SQL语句。

```sql
    --查询比B班住在东京的所有学生年龄都小的A班学生
    SELECT ＊
      FROM Class_A
     WHERE age < ALL ( SELECT age
                        FROM Class_B
                        WHERE city =’东京’);
```

执行结果

```sql
    name   age   city
    -----  ----  ----
    拉里    19     埼玉
```

如果山田的年龄仍是NULL时的分析

- ALL谓词其实是多个以AND连接的逻辑表达式的省略写法。

```sql
    --1．执行子查询获取年龄列表
    SELECT ＊
      FROM Class_A
     WHERE age < ALL ( 22, 23, NULL );

    --2．将ALL谓词等价改写为AND
    SELECT ＊
      FROM Class_A
     WHERE (age < 22) AND (age < 23) AND (age < NULL);

    --3．对NULL使用“<”后，结果变为 unknown
    SELECT ＊
      FROM Class_A
     WHERE (age < 22)  AND (age < 23) AND unknown;

    --4. 如果AND运算里包含unknown，则结果不为true
    SELECT ＊
      FROM Class_A
     WHERE false 或 unknown;
```

## 5. 限定谓词和极值函数不是等价的

如果用极值函数重写刚才的SQL

```sql
    --查询比B班住在东京的年龄最小的学生还要小的A班学生
    SELECT ＊
      FROM Class_A
     WHERE age < ( SELECT MIN(age)
                    FROM Class_B
                    WHERE city =’东京’);
```

执行结果

```sql
    name   age   city
    -----  ----  ----
    拉里    19    埼玉
    伯杰    21    千叶
```

没有问题。即使山田的年龄无法确定，这段代码也能查询到拉里和伯杰两人。这是因为，**极值函数在统计时会把为NULL的数据排除掉。**

ALL谓词和极值函数表达的命题含义

- ALL谓词：他的年龄比在东京住的所有学生都小——Q1
- 极值函数：他的年龄比在东京住的年龄最小的学生还要小——Q2

还有一种情况下它们也是不等价的

- 谓词（或者函数）的输入为空集的情况

Class_B没有住在东京的学生！

![](https://res.weread.qq.com/wrepub/epub_26211874_80)

- 使用ALL谓词的SQL语句会查询到A班的所有学生
- 然而用极值函数查询时一行数据都查询不到。

```sql
    --1．极值函数返回NULL
    SELECT ＊
      FROM Class_A
     WHERE age < NULL;

    --2．对NULL使用“<”后结果为 unknown
    SELECT ＊
      FROM Class_A
     WHERE unknown;
```

## 6. 聚合函数和NULL

COUNT以外的聚合函数在输入为空表时都返回NULL

```sql
    --查询比住在东京的学生的平均年龄还要小的A班学生的SQL语句？
    SELECT ＊
      FROM Class_A
     WHERE age < ( SELECT AVG(age)
                    FROM Class_B
                    WHERE city =’东京’);
```

# 本节小结
