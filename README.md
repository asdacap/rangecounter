Range Tree Counter
==================

Its an experiment to create a module that simply counts event, and enable query for how many event occur in an interval
of time. A basic implementation that I saw is to use a bucket per unit time (such as seconds or minute) to count the 
event. And to query for an interval, we fetch all bucket within the range of time. I thought, this looks inefficient,
if we need to query for 60 minute, we need to fetch 60 records. So I thought of using a [segment tree](https://en.wikipedia.org/wiki/Segment_tree), (which I mistakenly thought was called a range tree)
and here is a basic implementation. This is not the only solution though, simply using a relational database that
is properly indexed can yield a O(log n) complexity also. But it may be write limited as relational database tends to.

Why?
----
Because sometimes you just need to code it out of your head.

Use case
--------

So the use case that I'm looking at have a relatively small interval, which is why using a simple bucket is not a problem.
The query generally something along the line of how many event do this customer have for the past 7 days, and in that
case, the bucket is per day, and it will only need to fetch 7 query. It also use redis pipeline, so the latency is 
negligible. In addition to that, a write to a segment tree is O(log n), while this bucket scheme only need a single
write operation. So it is unclear if a redis backed segment tree will bring any benefit, so that is why I made this.

Segment Tree
------------

A big downside of a segment tree is the O(log n) write. This can be mitigated somehow by limiting the height of the tree
and increasing the number of child node per node. This is similar to how b-tree is used on relational database instead
of binary tree. Well b-tree is actually used for another unrelated reason, but it have the effect of lowering the height
of the tree. 

Limiting the height of the segment tree means that if the query interval is really large, it will be kinda the same as the 
bucket scheme. In theory, if the max height is 1, it will perform the same as the bucket scheme.

Increasing the number of child node per node should increase the number of node required to read a particular interval. 
But as we shall see, on random test, sometimes, it perform the same. It theory, if it is too high, it would again,
be more or less like the bucket scheme.

Another thing to consider is the number of key or node used. In a bucket scheme, a bucket only use a key. In a segment
tree, it adds a O(log n) factor to it. How much exactly, depends on workload.

It is also possible to reuse an existing tree for a different time unit. For example, using a seconds tree to query 
an hour. Since a segment tree will run in O(log n), it should not be too slow.

Benchmark
---------

The benchmark is run using the conventional Go benchmark utility. It use the new ReportMetric function which is only 
available starting from Go 1.13. 

Each benchmark is 10000 round of random increment and 10000 round of random query. The per-bucket unit of time is minutes.
But there are some configuration which query by the minutes, but stores by the seconds, adding 60 factor to the query 
interval. Three query interval is tested, 5, 20 and 100. Query interval here refers to the size of the query, as in, how many 
minutes to look for the result. 5 minutes, 20 minutes, and 100 minutes. Really, the unit of time here does not matter except 
for comparison when using different underlying store. Additionally these are maximum query interval. The benchmark randomly
calls multiple interval up to the specified maximum query interval. The benchmark reports the number of key used during increment 
and query. It also report the number of total keys used by the store. These number are per-query/increment number as they
have been divided by 10000, so its average number. 

The tree implementation are named in (h-p) format, where h is the tree height and p is the child key length, which determine
the number of child per node which is 2^p.

Implementation (Read/Write/KeyUsed) |               5 |              20 |             100
------------------------------------|----------------:|----------------:|-----------------:
Bucket                              | 2.98/1.00/0.951 | 10.4/1.00/0.951 | 50.1/1.00/0.951
Tree (2-1)                          | 2.68/2.00/1.86  |  6.64/2.00/1.86 |  26.5/2.00/1.86
Tree (2-2)                          | 2.98/2.00/1.77  |  6.00/2.00/1.77 |  16.2/2.00/1.77
Tree (4-1)                          | 2.68/4.00/3.37  |  4.76/4.00/3.37 |  10.1/4.00/3.37
Tree (8-1)                          | 2.68/8.00/4.40  |  4.74/8.00/4.40 |  7.12/8.00/4.40
Tree (16-1)                         | 2.68/16.0/4.48  |  4.74/16.0/4.48 |  7.12/16.0/4.48
Tree (8-2)                          | 2.98/8.00/2.48  |  5.95/8.00/2.48 |  9.49/8.00/2.48
Tree (4-4)                          | 2.98/4.00/1.49  |  10.2/4.00/1.49 |  17.9/4.00/1.49
Bucket (to seconds)                 | 179/1.00/0.951  |  626/1.00/0.951 | 3005/1.00/0.951
Tree (8-1) (to seconds)             | 11.2/8.00/7.55  |  14.8/8.00/7.55 |  33.5/8.00/7.55
Tree (16-1) (to seconds)            | 11.2/16.0/10.1  |  12.9/16.0/10.1 |  15.1/16.0/10.1
Tree (32-1) (to seconds)            | 11.2/32.0/10.1  |  12.9/32.0/10.1 |  15.1/32.0/10.1
Tree (8-2) (to seconds)             | 15.6/8.00/5.27  |  18.1/8.00/5.27 |  21.4/8.00/5.27
Tree (16-2) (to seconds)            | 15.6/16.0/5.29  |  18.1/16.0/5.29 |  21.4/16.0/5.29
Tree (4-4) (to seconds)             | 29.6/4.00/2.86  |  35.3/4.00/2.86 |  44.5/4.00/2.86
Tree (8-4) (to seconds)             | 29.6/8.00/2.87  |  35.3/8.00/2.87 |  44.5/8.00/2.87
Tree (4-8 (to seconds)              |  171/4.00/1.77  |   241/4.00/1.77 |   267/4.00/1.77

Read wise, we can see that even with tree height 2, there is about 10% improvement. But write increase by a factor of 2.
Increasing the tree height does not help much at all, but they significantly increase the write time.

On higher query interval, we starts to see the benefit of the segment tree. Both tree with height of two reduce read by about
35%. Interestingly the tree 4-1, 8-1, 16-1 are more or less identical. 2^4 is 16, while 2^5 is 32. This means when query
interval is 20 (or on average here 10), it won't benefit much from a tree with higher height than 4. (8-2), changes this
calculation a bit to (2^2)^2, which is 16. Which means, its only using 2 level of the tree instead of 4. For the tree
(4-4), its only using one level. Write wise, the trees write remain consistent throughout the query interval. One thing to 
note is the number of key used. Pretty much all tree config uses more keys than the bucket. Smaller child key size
increase the number of key used too. If read is a major concern instead, and the number of key used is not a concern,
and you know your interval won't be more than 20, or so, the the tree (4-1) would be a good idea.

On query interval of 100 we could really see the benefit of using a segment tree. Even using only height 2, the read touches
twice less keys than the bucket. Tree height 8+ with child key size 1 shows the best read performance. However it increase
write factor by 7 fold. Interestingly increasing the child key size actually reduce the read performance. This is because
for each level, it may need to read up to an additional 3 key instead of only 1 for when key size is 1. This is variable
though as it depends if the query interval align or not. If the query interval align perfectly, it may need to read exactly one
key. This segment alignment is one reason why I use a random test. I've always wondered how much does the alignment degrades
performance. It turns out, it is observable. That said, it would have twice lower writes than read. 16 vs 8 writes is 
pretty significant and could be worth the read downgrade. That said, at this interval, even (4-4) improves upon basic bucket.
In fact, (2-4) would performs the same here, and only doubles the write.

As for the tree with seconds as underlying store unit, it depends on the range. (16-1) seems to perform the best. The
difference between this and just setting the max query interval to say.. 300 is that the minimum query interval is 60. Also
I wonder if the minute to second alignment affect much. It turns out, it definitely degrades read over writes. On 
extreme case, the (4-8), it is 5 times worst than (4-4). Making (4-4) much better option in this case. At least it is
better that the raw bucket scheme, which is 8 times slower than (4-8) on higher interval, but perform the same on lower
interval.

Bottomline
----------

Its not a silver bullet. Highly depends on your workload. If you have a very variable query interval and writes is not a
bit problem, then it may worth the trouble. If writes and number of keys is not a problem, then a segment tree of child
key size 1 will on average improve read performance. The height of the tree, depends on the interval. If you set it more 
than you need, then its a waste of money, and writes. Its all an act of balance between writes, key requirement, reads,
tree height and child key size.
