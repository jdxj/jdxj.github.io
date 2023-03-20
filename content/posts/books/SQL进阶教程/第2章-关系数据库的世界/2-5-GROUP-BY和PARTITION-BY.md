---
title: "2-5 GROUP BY和PARTITION BY.md"
date: 2023-03-01T21:54:29+08:00
---

本节将以集合论和群论中的“类”这一重要概念为核心，阐明GROUP BY和PARTITION BY的意义。

SQL的语句中具有分组功能的是GROUP BY和PARTITION BY，它们都可以根据指定的列为表分组。区别仅仅在于，GROUP BY在分组之后会把每个分组聚合成一
行数据。

Teams

![](https://res.weread.qq.com/wrepub/epub_26211874_351)

```sql
    SELECT member, team, age ,
          RANK() OVER(PARTITION BY team ORDER BY age DESC) rn,
          DENSE_RANK() OVER(PARTITION BY team ORDER BY age DESC) dense_rn,
          ROW_NUMBER() OVER(PARTITION BY team ORDER BY age DESC) row_num
      FROM Members
     ORDER BY team, rn;
```

执行结果

![](https://res.weread.qq.com/wrepub/epub_26211874_352)

分割后的子集

![](https://res.weread.qq.com/wrepub/epub_26211874_353)

重点关注一下划分出的子集

1. 它们全都是非空集合。
2. 所有子集的并集等于划分之前的集合。
3. 任何两个子集之间都没有交集。

对3取余给自然数集合N分类

- 余0的类：M1 = {0, 3, 6, 9, …}
- 余1的类：M2 = {1, 4, 7, 10, …}
- 余2的类：M2 = {2, 5, 8, 11, …}

从类的第2个性质我们知道，这3个类涵盖了全部自然数。

- M1 + M2 + M3 = N

MOD函数

```sql
    --对从1到10的整数以3为模求剩余类
    SELECT MOD(num, 3) AS modulo,
          num
      FROM Natural
     ORDER BY modulo, num;
```

执行结果

![](https://res.weread.qq.com/wrepub/epub_26211874_355)

随机地将数据减为原来的五分之一

```sql
    --从原来的表中抽出（大约）五分之一行的数据
    SELECT ＊
      FROM SomeTbl
     WHERE MOD(seq, 5) = 0;


    --表中没有连续编号的列时，使用ROW_NUMBER函数就可以了
    SELECT ＊
      FROM (SELECT col,
                  ROW_NUMBER() OVER(ORDER BY col) AS seq
              FROM SomeTbl)
     WHERE MOD(seq, 5) = 0;
```
