---
title: "1-8 EXISTS谓词的用法"
date: 2023-03-06T21:42:02+08:00
draft: false
---

# 理论篇

## 什么是谓词

谓词是一种特殊的函数，返回值是真值。

- true
- false
- unknown

在关系数据库里，表中的一行数据可以看作是一个命题。

Tbl_A

![](https://res.weread.qq.com/wrepub/epub_26211874_205)

## 实体的阶层

同样是谓词，但是与=、BETWEEN等相比，EXISTS的用法还是大不相同的。概括来说，区别在于“谓词的参数可以取什么值”。

```sql
    SELECT id
      FROM Foo F
     WHERE EXISTS
            (SELECT ＊
              FROM Bar B
              WHERE F.id=B.id );
```

看一下EXISTS()的括号中的内容

```sql
              SELECT ＊
                FROM Bar B
               WHERE A.id = T2.id
```

在EXISTS的子查询里，SELECT子句的列表可以有下面这三种写法。

- 通配符：SELECT ＊
- 常量：SELECT ‘这里的内容任意’
- 列名：SELECT col

不管采用上面这三种写法中的哪一种，得到的结果都是一样的。

![](https://res.weread.qq.com/wrepub/epub_26211874_207)

=或者BETWEEEN等输入值为一行的谓词叫作“一阶谓词”，而像EXISTS这样输入值为行的集合的谓词叫作“二阶谓词”。阶（order）是用来区分集合或谓词的阶
数的概念。

- 三阶谓词＝输入值为“集合的集合”的谓词
- 四阶谓词＝输入值为“集合的集合的集合”的谓词
- ……

我们可以像上面这样无限地扩展阶数，但是SQL里并不会出现三阶以上的情况

EXISTS因接受的参数是集合这样的一阶实体而被称为二阶谓词，但是谓词也是函数的一种，因此我们也可以说EXISTS是高阶函数。

![](https://res.weread.qq.com/wrepub/epub_26211874_208)

## 全称量化和存在量化

“所有的x都满足条件P”或者“存在（至少一个）满足条件P的x”。前者称为“全称量词”，后者称为“存在量词”，分别记作∀、∃。这两个符号看起来很奇怪。其
实，全称量词的符号其实是将字母A上下颠倒而形成的，存在量词则是将字母E左右颠倒而形成的。

sql没有实现全称量词, 但是没有全称量词并不算是SQL的致命缺陷。因为全称量词和存在量词只要定义了一个，另一个就可以被推导出来。

- ∀ xPx = ¬ ∃ x¬P（所有的x都满足条件P＝不存在不满足条件P的x）
- ∃ xPx = ¬ ∀ x¬Px（存在x满足条件P＝并非所有的x都不满足条件P）

# 实践篇

## 查询表中“不”存在的数据

在有些情况下，我们不得不从表中查找出“不存在的数据”。

Meetings

![](https://res.weread.qq.com/wrepub/epub_26211874_212)

目标结果如下所示

```sql
    meeting          person
    ----------       --------
    第1次             宫田
    第2次             坂东
    第2次             水岛
    第3次             伊藤
```

思路是先假设所有人都参加了全部会议，并以此生成一个集合，然后从中减去实际参加会议的人。这样就能得到缺席会议的人。

所有人都参加了全部会议的集合可以通过下面这样的交叉连接来求得。

```sql
    SELECT DISTINCT M1.meeting, M2.person
      FROM Meetings M1 CROSS JOIN Meetings M2;
```

所有人都参加了全部会议时

![](https://res.weread.qq.com/wrepub/epub_26211874_213)

然后我们从这张表中减掉实际参会者的集合

```sql
    --求出缺席者的SQL语句(1)：存在量化的应用
    SELECT DISTINCT M1.meeting, M2.person
      FROM Meetings M1 CROSS JOIN Meetings M2
     WHERE NOT EXISTS
          (SELECT ＊
              FROM Meetings M3
            WHERE M1.meeting = M3.meeting
              AND M2.person = M3.person);
```

这道例题还可以用集合论的方法来解答，即像下面这样使用差集运算。

```sql
    ----求出缺席者的SQL语句(2)：使用差集运算
    SELECT M1.meeting, M2.person
      FROM Meetings M1, Meetings M2
    EXCEPT
    SELECT meeting, person
      FROM Meetings;
```

## 全称量化(1)：习惯“肯定⇔双重否定”之间的转换

学生考试成绩的表

TestScores

![](https://res.weread.qq.com/wrepub/epub_26211874_214)

请查询出“所有科目分数都在50分以上的学生”

将查询条件“所有科目分数都在50分以上”转换成它的双重否定“没有一个科目分数不满50分”，然后用NOT EXISTS来表示转换后的命题。

```sql
    SELECT DISTINCT student_id
      FROM TestScores TS1
     WHERE NOT EXISTS                --不存在满足以下条件的行
            (SELECT ＊
              FROM TestScores TS2
              WHERE TS2.student_id = TS1.student_id
                AND TS2.score < 50);    --分数不满50分的科目
```

执行结果

```sql
    student_id
    -----------
          100
          200
          400
```

查询出满足下列条件的学生。

- 数学的分数在80分以上。
- 语文的分数在50分以上。

针对同一个集合内的行数据进行了条件分支后的全称量化。

- “某个学生的所有行数据中，如果科目是数学，则分数在80分以上；如果科目是语文，则分数在50分以上。”

```sql
    CASE WHEN subject =’数学’AND score >= 80 THEN 1
        WHEN subject =’语文’AND score >= 50 THEN 1
        ELSE 0 END
```

首先，数学和语文之外的科目不在我们考虑范围之内，所以通过IN条件进行一下过滤。然后，通过子查询来描述“数学80分以上，语文50分以上”这个条件。

```sql
    SELECT DISTINCT student_id
      FROM TestScores TS1
     WHERE subject IN (’数学’, ’语文’)
      AND NOT EXISTS
            (SELECT ＊
              FROM TestScores TS2
              WHERE TS2.student_id = TS1.student_id
                AND 1 = CASE WHEN subject =’数学’AND score < 80 THEN 1
                            WHEN subject =’语文’AND score < 50 THEN 1
                            ELSE 0 END);
```

排除掉没有语文分数的学号为400的学生。

```sql
    SELECT student_id
      FROM TestScores TS1
     WHERE subject IN (’数学’, ’语文’)
      AND NOT EXISTS
            (SELECT ＊
              FROM TestScores TS2
              WHERE TS2.student_id = TS1.student_id
                AND 1 = CASE WHEN subject =’数学’AND score < 80 THEN 1
                            WHEN subject =’语文’AND score < 50 THEN 1
                            ELSE 0 END)
     GROUP BY student_id
    HAVING COUNT(＊) = 2;   --必须两门科目都有分数


    student_id
    ----------
          100
          200
```

## 全称量化(2)：集合VS谓词——哪个更强大？

项目工程管理表

Projects

![](https://res.weread.qq.com/wrepub/epub_26211874_216)

从这张表中查询出哪些项目已经完成到了工程1

Joe Celko曾经借助HAVING子句用面向集合的方法进行过解答

```sql
    --查询完成到了工程1的项目：面向集合的解法
    SELECT project_id
      FROM Projects
     GROUP BY project_id
    HAVING COUNT(＊) = SUM(CASE WHEN step_nbr <= 1 AND status =’完成’THEN 1
                            WHEN step_nbr  > 1 AND status =’等待’THEN 1
                            ELSE 0 END);
```

执行结果

```sql
    project_id
    -----------
    CS300
```

针对每个项目，将工程编号为1以下且状态为“完成”的行数，和工程编号大于1且状态为“等待”的行数加在一起，如果和等于该项目数据的总行数，则该项目符合
查询条件。

用谓词逻辑

```sql
    --查询完成到了工程1的项目：谓词逻辑的解法
    SELECT ＊
      FROM Projects P1
     WHERE NOT EXISTS
          (SELECT status
            FROM Projects P2
            WHERE P1.project_id = P2. project_id      --以项目为单位进行条件判断
              AND status <> CASE WHEN step_nbr <= 1   --使用双重否定来表达全称量化命题
                              THEN ’完成’
                              ELSE ’等待’ END);
```

执行结果

```sql
    project_id    step_nbr    status
    -----------   --------    ------
    CS300                  0    完成
    CS300                  1    完成

    CS300                  2    等待
    CS300                  3    等待
```

## 对列进行量化：查询全是1的行

ArrayTbl

![](https://res.weread.qq.com/wrepub/epub_26211874_217)

在使用这种模拟数组的表时遇到的需求一般都是下面这两种形式。

- 查询“都是1”的行。
- 查询“至少有一个9”的行。

```sql
    --“列方向”的全称量化：不优雅的解答
    SELECT ＊
      FROM ArrayTbl
     WHERE col1 = 1
      AND col2 = 1
          ·
          ·
          ·
      AND col10 = 1;
```

SQL语言其实还准备了一个谓词，帮助我们进行“列方向”的量化。

```sql
    --“列方向”的全称量化：优雅的解答
    SELECT ＊
      FROM ArrayTbl
     WHERE 1 = ALL (col1, col2, col3, col4, col5, col6, col7, col8, col9, col10);
```

```sql
    key  col1  col2  col3  col4  col5  col6  col7  col8  col9  col10
    ---  ----  ----  ----  ----  ----  ----  ----  ----  ----  -----
      C     1     1     1     1     1     1     1     1     1      1
```

如果想表达“至少有一个9”这样的存在量化命题，可以使用ALL的反义谓词ANY。

```sql
    --列方向的存在量化(1)
    SELECT ＊
      FROM ArrayTbl
     WHERE 9 = ANY (col1, col2, col3, col4, col5, col6, col7, col8, col9, col10);
```

```sql
    key  col1  col2  col3  col4  col5  col6  col7  col8  col9  col10
    ---  ----  ----  ----  ----  ----  ----  ----  ----  ----  -----
      D                 9
      E           3           1     9                 9
```

或者也可以使用IN谓词代替ANY。

```sql
    --列方向的存在量化(2)
    SELECT ＊
      FROM ArrayTbl
     WHERE 9 IN (col1, col2, col3, col4, col5, col6, col7, col8, col9, col10);
```

如果左边不是具体值而是NULL，这种写法就不行了。

```sql
    --查询全是NULL的行：错误的解法
    SELECT ＊
      FROM ArrayTbl
     WHERE NULL = ALL (col1, col2, col3, col4, col5, col6, col7, col8, col9, col10);
```

不管表里的数据是什么样的，这条SQL语句的查询结果都是空。这是因为，ALL谓词会被解释成col1 = NULL AND col2 = NULL AND ……col10 = NULL。
这种情况下，我们需要使用COALESCE函数。

```sql
    --查询全是NULL的行：正确的解法
    SELECT ＊
      FROM ArrayTbl
     WHERE COALESCE(col1, col2, col3, col4, col5, col6, col7, col8, col9, col10) IS NULL;
```

```sql
    key  col1  col2  col3  col4  col5  col6  col7  col8  col9  col10
    ---  ----  ----  ----  ----  ----  ----  ----  ----  ----  -----
     A
```

# 本节小结

- SQL中的谓词指的是返回真值的函数。
- EXISTS与其他谓词不同，接受的参数是集合。
- 因此EXISTS可以看成是一种高阶函数。
- SQL中没有与全称量词相当的谓词，可以使用NOT EXISTS代替。
