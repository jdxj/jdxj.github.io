---
title: "1-9 用SQL处理数列"
date: 2023-03-07T21:49:48+08:00
draft: false
---

# 生成连续编号

谜题：00～99的100个数中，0, 1, 2, …, 9这10个数字分别出现了多少次？

00～99的数中，数字0～9各出现了20次

![](https://res.weread.qq.com/wrepub/epub_26211874_230)

无论多大的数，都可以由这张表中的10个数字组合而成。

Digits

![](https://res.weread.qq.com/wrepub/epub_26211874_231)

通过对两个Digits集合求笛卡儿积而得出0～99的数字

```sql
    --求连续编号(1)：求0~99的数
    SELECT D1.digit + (D2.digit ＊ 10)  AS seq
      FROM Digits D1 CROSS JOIN Digits D2
     ORDER BY seq;
```

执行结果

```sql
    seq
    ---
      0
      1
      2
     ：
     ：
     ：
     98
     99
```

笛卡儿积：得到所有可能的组合

![](https://res.weread.qq.com/wrepub/epub_26211874_232)

如果只想生成从1开始，或者到542结束的数，只需在WHERE子句中加入过滤条件就可以了。

```sql
    --求连续编号(2)：求1~542的数
    SELECT D1.digit + (D2.digit ＊ 10) + (D3.digit ＊ 100) AS seq
      FROM Digits D1 CROSS JOIN Digits D2
            CROSS JOIN Digits D3
     WHERE D1.digit + (D2.digit ＊ 10)
                    + (D3.digit ＊ 100) BETWEEN 1 AND 542
     ORDER BY seq;
```

通过将这个查询的结果存储在视图里，就可以在需要连续编号时通过简单的SELECT来获取需要的编号。

```sql
--生成序列视图（包含0~999）
CREATE VIEW Sequence (seq)
AS SELECT D1.digit + (D2.digit ＊ 10) + (D3.digit ＊ 100)
    FROM Digits D1 CROSS JOIN Digits D2
            CROSS JOIN Digits D3;

--从序列视图中获取1~100
SELECT seq
FROM Sequence
WHERE seq BETWEEN 1 AND 100
ORDER BY seq;
```

# 求全部的缺失编号

假设存在下面这样一张编号有缺失的表。

Seqtbl

![](https://res.weread.qq.com/wrepub/epub_26211874_235)

利用序列视图

```sql
--EXCEPT版
SELECT seq
  FROM Sequence
 WHERE seq BETWEEN 1 AND 12
EXCEPT
SELECT seq FROM SeqTbl;

--NOT IN版
SELECT seq
FROM Sequence
WHERE seq BETWEEN 1 AND 12
  AND seq NOT IN (SELECT seq FROM SeqTbl);
```

执行结果

```sql
    seq
    ---
      3
      9
     10
```

可能像下面这么做性能会有所下降，但是通过扩展BETWEEN谓词的参数，我们可以动态地指定目标表的最大值和最小值。

```sql
    --动态地指定连续编号范围的SQL语句
    SELECT seq
      FROM Sequence
     WHERE seq BETWEEN (SELECT MIN(seq) FROM SeqTbl)
                  AND (SELECT MAX(seq) FROM SeqTbl)
    EXCEPT
    SELECT seq FROM SeqTbl;
```

- 这种写法在查询上限和下限未必固定的表时非常方便。两个自查询没有相关性，而且只会执行一次。
- 如果在“seq”列上建立索引，那么极值函数的运行可以变得更快速。

# 三个人能坐得下吗

火车座位预订情况的表

Seats

![](https://res.weread.qq.com/wrepub/epub_26211874_237)

问题是，从1～15的座位编号中，找出连续3个空位的全部组合。

希望得到的结果是

- 3～5
- 7～9
- 8～10
- 9～11

假设所有的座位排成了一条直线。

7～11的序列包含3个子序列

![](https://res.weread.qq.com/wrepub/epub_26211874_238)

借助上面的图表我们可以知道，需要满足的条件是，以n为起点、n+(3-1)为终点的座位全部都是未预订状态

```sql
    --找出需要的空位(1)：不考虑座位的换排
    SELECT S1.seat   AS start_seat, '~', S2.seat AS end_seat
      FROM Seats S1, Seats S2
     WHERE S2.seat = S1.seat + (:head_cnt -1)  --决定起点和终点
      AND NOT EXISTS
              (SELECT ＊
                FROM Seats S3
                WHERE S3.seat BETWEEN S1.seat AND S2.seat
                  AND S3.status <>’未预订’);
```

发生换排的情况。假设这列火车每一排有5个座位。我们在表中加上表示行编号“row_id”列。

Seats2

![](https://res.weread.qq.com/wrepub/epub_26211874_239)

因为发生换排，所以9~11的序列不符合条件

![](https://res.weread.qq.com/wrepub/epub_26211874_240)

```sql
    --找出需要的空位(2)：考虑座位的换排
SELECT S1.seat   AS start_seat, '~', S2.seat AS end_seat
FROM Seats2 S1, Seats2 S2
WHERE S2.seat = S1.seat + (:head_cnt -1)  --决定起点和终点
  AND NOT EXISTS
    (SELECT ＊
     FROM Seats2 S3
     WHERE S3.seat BETWEEN S1.seat AND S2.seat
       AND (    S3.status <>’未预订’
         OR S3.row_id <> S1.row_id));
```

执行结果

```sql
    start_seat '~'  end_seat
    ----------  ---  --------
            3  ~     5
            8  ~    10
            11  ~    13
```

序列内的点需要满足的条件是，“所有座位的状态都是‘未预订’，且行编号相同”。这里新加的条件是“行编号相同”，等价于“与起点的行编号相同”（当然，与终
点的行编号相同也可以）。把这个条件直接写成SQL语句的话，就是像下面这样。

```sql
    S3.status =’未预订’AND S3.row_id = S1.row_id
```

由于SQL中不存在全称量词，所以我们必须使用这个条件的否定，即改成下面这样的否定形式。

```sql
    NOT (S3.status =’未预订’AND S3.row_id = S1.row_id)
  = S3.status <>’未预订’OR S3.row_id <> S1.row_id
```

# 最多能坐下多少人

这次要查询的是“按现在的空位状况，最多能坐下多少人”。换句话说，要求的是最长的序列。

Seats3

![](https://res.weread.qq.com/wrepub/epub_26211874_242)

要想保证从座位A到另一个座位B是一个序列，则下面的3个条件必须全部都满足。

- 条件1：起点到终点之间的所有座位状态都是“未预订”。
- 条件2：起点之前的座位状态不是“未预订”。
- 条件3：终点之后的座位状态不是“未预订”。

不符合条件1的情况

![](https://res.weread.qq.com/wrepub/epub_26211874_243)

不符合条件2的情况

![](https://res.weread.qq.com/wrepub/epub_26211874_244)

不符合条件3的情况

![](https://res.weread.qq.com/wrepub/epub_26211874_245)

先生成一张下面这样的视图。

```sql
    --第一阶段：生成存储了所有序列的视图
    CREATE VIEW Sequences (start_seat, end_seat, seat_cnt) AS
    SELECT S1.seat  AS start_seat,
          S2.seat  AS end_seat,
          S2.seat - S1.seat + 1 AS seat_cnt
      FROM Seats3 S1, Seats3 S2
     WHERE S1.seat <= S2.seat  --第一步：生成起点和终点的组合

        AND NOT EXISTS    --第二步：描述序列内所有点需要满足的条件
            (SELECT ＊
              FROM Seats3 S3
              WHERE (     S3.seat BETWEEN S1.seat AND S2.seat
                      AND S3.status <>’未预订’)  --条件1的否定
                OR  (S3.seat = S2.seat + 1 AND S3.status =’未预订’)
                                                        --条件3的否定
                OR  (S3.seat = S1.seat -1 AND S3.status =’未预订’));
                                                        --条件2的否定
```

这个视图包含以下的内容。

```sql
    start_seat    end_seat    seat_cnt
    ------------  ----------  ----------
            2            5            4
            7            7            1
            9           10            2
```

我们从这个视图中找出座位数（seat_cnt）最大的一行数据。

```sql
    --第二阶段：求最长的序列
    SELECT start_seat, '~', end_seat, seat_cnt
      FROM Sequences
     WHERE seat_cnt = (SELECT MAX(seat_cnt) FROM Sequences);
```

# 单调递增和单调递减

某公司股价动态的表

MyStock

![](https://res.weread.qq.com/wrepub/epub_26211874_247)

求一下股价单调递增的时间区间。从上表来看，目标结果是下面两个。

- 2007-01-06～2007-01-08
- 2007-01-14～2007-01-17

首先进行第一步——通过自连接生成起点和终点的组合。

```sql
    --生成起点和终点的组合的SQL语句
    SELECT S1.deal_date  AS start_date,
          S2.deal_date  AS end_date
      FROM MyStock S1, MyStock S2
     WHERE S1.deal_date < S2.deal_date;
```

第二步——描述起点和终点之间的所有点需要满足的条件。

- 对于区间内的任意两个时间点，命题“较晚时间的股价高于较早时间的股价”都成立。
- 然后，我们将这个条件反过来，得到需要的条件——区间内不存在两个时间点使得较早时间的股价高于较晚时间的股价。

```sql
    --求单调递增的区间的SQL语句：子集也输出
    SELECT S1.deal_date   AS start_date,
          S2.deal_date   AS end_date
      FROM MyStock S1, MyStock S2
     WHERE S1.deal_date < S2.deal_date  --第一步：生成起点和终点的组合
      AND  NOT EXISTS
              ( SELECT ＊  --第二步：描述区间内所有日期需要满足的条件
                  FROM MyStock S3, MyStock S4
                  WHERE S3.deal_date BETWEEN S1.deal_date AND S2.deal_date
                  AND S4.deal_date BETWEEN S1.deal_date AND S2.deal_date
                    AND S3.deal_date < S4.deal_date
                    AND S3.price >= S4.price);
```

执行结果

```sql
    start_date     end_date
    ------------   -------------
    2007-01-06     2007-01-08
    2007-01-14     2007-01-16
    2007-01-14     2007-01-17
    2007-01-16     2007-01-17
```

最后，我们要把这些不需要的子集排除掉。使用极值函数很容易就能实现。

```sql
    --排除掉子集，只取最长的时间区间
    SELECT MIN(start_date) AS start_date,      --最大限度地向前延伸起点
          end_date
      FROM  (SELECT S1.deal_date AS start_date,
                    MAX(S2.deal_date) AS end_date  --最大限度地向后延伸终点
              FROM MyStock S1, MyStock S2
              WHERE S1.deal_date < S2.deal_date
                AND NOT EXISTS
                (SELECT ＊
                    FROM MyStock S3, MyStock S4
                  WHERE S3.deal_date BETWEEN S1.deal_date AND S2.deal_date
                    AND S4.deal_date BETWEEN S1.deal_date AND S2.deal_date
                    AND S3.deal_date < S4.deal_date
                    AND S3.price >= S4.price)
            GROUP BY S1.deal_date) TMP
    GROUP BY end_date;
```

执行结果

```sql
    start_date     end_date
    ------------   -------------
    2007-01-06     2007-01-08
    2007-01-14     2007-01-17
```

# 本节小结

SQL处理数据的方法有两种。

- 第一种是把数据看成忽略了顺序的集合。
- 第二种是把数据看成有序的集合，此时的基本方法如下。
  - 首先用自连接生成起点和终点的组合
  - 其次在子查询中描述内部的各个元素之间必须满足的关系
- 要在SQL中表达全称量化时，需要将全称量化命题转换成存在量化命题的否定形式，并使用NOT EXISTS谓词。这是因为SQL只实现了谓词逻辑中的存在量词。
