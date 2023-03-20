---
title: "1-10 HAVING子句又回来了"
date: 2023-03-08T21:39:31+08:00
---

# 各队，全体点名

你需要做的是查出现在可以出勤的队伍。可以出勤即队伍里所有队员都处于“待命”状态。

Teams

![](https://res.weread.qq.com/wrepub/epub_26211874_254)

“所有队员都处于‘待命’状态”这个条件是全称量化命题，所以可以用NOT EXISTS来表达。

```sql
    -- 用谓词表达全称量化命题
    SELECT team_id, member
      FROM Teams T1
     WHERE NOT EXISTS
            (SELECT ＊
              FROM Teams T2
              WHERE T1.team_id = T2.team_id
                AND status <>’待命’);
```

执行结果

```sql
    team_id   member
    -------   ------
          3   简
          3   哈特
          3   迪克
          4   贝斯
```

使用HAVING子句

```sql
    -- 用集合表达全称量化命题(1)
    SELECT team_id
      FROM Teams
     GROUP BY team_id
    HAVING COUNT(＊) = SUM(CASE WHEN status =’待命’
                              THEN 1
                              ELSE 0 END);
```

执行结果

```sql
    team_id
    -------
          3
          4
```

第一步还是使用GROUP BY子句将Teams集合以队伍为单位划分成几个子集。

![](https://res.weread.qq.com/wrepub/epub_26211874_256)
![](https://res.weread.qq.com/wrepub/epub_26211874_257)

HAVING子句中的条件还可以像下面这样写。

```sql
    -- 用集合表达全称量化命题(2)
    SELECT team_id
      FROM Teams
     GROUP BY team_id
    HAVING MAX(status) =’待命’
      AND MIN(status) =’待命’;
```

极值函数可以使用参数字段的索引，所以这种写法性能更好（当然本例中只有3种值，建立索引也并没有太大的意义）。

也可以把条件放在SELECT子句里，以列表形式显示出各个队伍是否所有队员都在待命

```sql
    -- 列表显示各个队伍是否所有队员都在待命
    SELECT team_id,
          CASE WHEN MAX(status) =’待命’AND MIN(status) =’待命’
                THEN ’全都在待命’
                ELSE’队长！人手不够’END AS status
      FROM Teams
     GROUP BY team_id;
```

```sql
    team_id   status
    -------   --------------------------
        1    队长！人手不够
        2    队长！人手不够
        3    全都在待命
        4    全都在待命
        5    队长！人手不够
```

# 单重集合与多重集合

关系数据库中的集合是允许重复数据存在的多重集合。与之相反，通常意义的集合论中的集合不允许数据重复，被称为“单重集合”（这是笔者自己造的词，并非
公认的术语）。

生产地的材料库存的表

Materials

![](https://res.weread.qq.com/wrepub/epub_26211874_259)

为了在各生产地之间调整重复的材料，我们需要调查出存在重复材料的生产地。

按生产地分组

![](https://res.weread.qq.com/wrepub/epub_26211874_260)

“排除掉重复元素后和排除掉重复元素前元素个数不相同”。如果不存在重复的元素，不管是否加上DISTINCT可选项，COUNT的结果都是相同的。

```sql
    -- 选中材料存在重复的生产地
    SELECT center
      FROM Materials
     GROUP BY center
    HAVING COUNT(material) <> COUNT(DISTINCT material);
```

执行结果

```sql
    center
    ------
    东京
    名古屋
```

显示哪种材料重复

```sql
    SELECT center,
          CASE WHEN COUNT(material) <> COUNT(DISTINCT material) THEN’存在重复’
              ELSE’不存在重复’END AS status
      FROM Materials
     GROUP BY center;

    center          status
    ----------     ----------
    大阪             不存在重复
    东京             存在重复
    福冈             不存在重复
    名古屋           存在重复
```

这个问题也可以通过将HAVING改写成EXISTS的方式来解决。

```sql
    --存在重复的集合：使用EXISTS
    SELECT center, material
      FROM Materials M1
     WHERE EXISTS
            (SELECT ＊
              FROM Materials M2
              WHERE M1.center = M2.center
                AND M1.receive_date <> M2.receive_date
                AND M1.material = M2.material);
```

```sql
    center  material
    ------- ---------
    东京     锌
    东京     锌
    名古屋    钛
    名古屋    钢
    名古屋    钢
    名古屋    钛
```

# 寻找缺失的编号：升级版

1-4节介绍过下面这样一条查询数列的缺失编号的查询语句

```sql
    -- 如果有查询结果，说明存在缺失的编号
    SELECT’存在缺失的编号’AS gap
      FROM SeqTbl
    HAVING COUNT(＊) <> MAX(seq);
```

这条SQL语句有一个前提条件，即数列的起始值必须是1。

放宽这个限制条件，思考一下不管数列的最小值是多少，都能用来判断该数列是否连续的SQL语句。

(1)不存在缺失编号（起始值＝1）

![](https://res.weread.qq.com/wrepub/epub_26211874_262)

(2)存在缺失编号（起始值＝1）

![](https://res.weread.qq.com/wrepub/epub_26211874_263)

(3)不存在缺失编号（起始值<>1）

![](https://res.weread.qq.com/wrepub/epub_26211874_264)

(4)存在缺失编号（起始值<>1）

![](https://res.weread.qq.com/wrepub/epub_26211874_265)

如果数列的最小值和最大值之间没有缺失的编号，它们之间包含的元素的个数应该是“最大值－最小值+1”。

```sql
    -- 如果有查询结果，说明存在缺失的编号：只调查数列的连续性
    SELECT’存在缺失的编号’  AS gap
      FROM SeqTbl
    HAVING COUNT(＊) <> MAX(seq) - MIN(seq) + 1 ;
```

不论是否存在缺失的编号，都想要返回结果

```sql
    -- 不论是否存在缺失的编号都返回一行结果
    SELECT CASE WHEN COUNT(＊) = 0
                THEN ’表为空’
                WHEN COUNT(＊) <> MAX(seq) - MIN(seq) + 1
                THEN ’存在缺失的编号’
                ELSE’连续’END AS gap
      FROM SeqTbl;
```

改进一下查找最小的缺失编号的SQL语句，去掉起始值必须是1的限制。

```sql
    -- 查找最小的缺失编号：表中没有1时返回1
    SELECT CASE WHEN COUNT(＊) = 0 OR MIN(seq) > 1  -- 最小值不是1时→返回1
              THEN 1

                ELSE (SELECT MIN(seq +1)           -- 最小值是1时→返回最小的缺失编号
                        FROM SeqTbl S1
                      WHERE NOT EXISTS
                          (SELECT ＊
                              FROM SeqTbl S2
                            WHERE S2.seq = S1.seq + 1))  END
        FROM SeqTbl;
```

这条SQL语句会返回下面这样的结果

- 情况(1):6（没有缺失的编号，所以返回最大值5的下一个数）
- 情况(2):3（最小的缺失编号）
- 情况(3):1（因为表中没有1）
- 情况(4):1（因为表中没有1）

# 为集合设置详细的条件

学生考试成绩的表

TestResults

![](https://res.weread.qq.com/wrepub/epub_26211874_267)

第1题：请查询出75%以上的学生分数都在80分以上的班级。

```sql
    SELECT class
      FROM TestResults
  GROUP BY class
    HAVING COUNT(＊) ＊ 0.75
          <= SUM(CASE WHEN score >= 80
                      THEN 1
                      ELSE 0 END) ;
```

执行结果

```sql
    class
    -------
        B
```

第2题：请查询出分数在50分以上的男生的人数比分数在50分以上的女生的人数多的班级。

```sql
    SELECT class
      FROM TestResults
  GROUP BY class
    HAVING SUM(CASE WHEN score >= 50 AND sex =’男’
                    THEN 1
                    ELSE 0 END)
        > SUM(CASE WHEN score >= 50 AND sex =’女’
                    THEN 1
                    ELSE 0 END) ;
```

执行结果

```sql
    class
    -------
        B
        C
```

第3题：请查询出女生平均分比男生平均分高的班级。

```sql
    -- 比较男生和女生平均分的SQL语句(1)：对空集使用AVG后返回0
      SELECT class
        FROM TestResults
    GROUP BY class
      HAVING AVG(CASE WHEN sex =’男’
                      THEN score
                      ELSE 0 END)
          < AVG(CASE WHEN sex =’女’
                      THEN score
                      ELSE 0 END) ;
```

执行结果

```sql
    class
    -------
        A
        D
```

根据标准SQL的定义，对空集使用AVG函数时，结果会返回NULL

```sql
    -- 比较男生和女生平均分的SQL语句(2)：对空集求平均值后返回NULL
      SELECT class
        FROM TestResults
    GROUP BY class
      HAVING AVG(CASE WHEN sex =’男’
                      THEN score
                      ELSE NULL END)
          < AVG(CASE WHEN sex =’女’
                      THEN score
                      ELSE NULL END) ;
```

这回D班男生的平均分是NULL。因此不管女生的平均分多少，D班都会被排除在查询结果之外。

# 本节小结

用于调查集合性质的常用条件及其用途

![](https://res.weread.qq.com/wrepub/epub_26211874_269)

在SQL中指定搜索条件时，最重要的是搞清楚搜索的实体是集合还是集合的元素。

- 如果一个实体对应着一行数据→那么就是元素，所以使用WHERE子句。
- 如果一个实体对应着多行数据→那么就是集合，所以使用HAVING子句。
- HAVING子句可以通过聚合函数（特别是极值函数）针对集合指定各种条件。
- 如果通过CASE表达式生成特征函数，那么无论多么复杂的条件都可以描述。
- HAVING子句很强大。
