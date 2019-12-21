#### crawl

爬虫

1,爬虫任务的统一调度

    engine统一调度所有资源,查询时间

2,爬虫任务的统一去重

    上一次爬去的内容留取一份在redis，用去下一次爬去任务时候做判断，是否为新数据，不同的新闻网可能会有重复，去除重复内容，由于每个网站的题目都是不一致的
    直接使用文章题目来判断是否是同一个内容比较难，

3,存储问题

    缓存保留该站点最近一次的内容，其他内容持久化到硬盘，包括存进数据库，保存成文件

4,断点续爬

    使用多线程，每个站点一个线程，线程间相互独立爬取，不影响内容

5,数据库保存格式

    时间 | 内容 | 来源


#### 结构解析：

由于golang对线程支持很好，不像python需要其他库才能做到多线程，所以爬虫模块就不需要下载和解析模块，每个任务都是一个goroutine即可

每次新加爬虫任务，则新加一个对应的网站逻辑，并且在一个列表结构中增加一个该爬虫函数的启动

每一站点的爬虫是单独的一个go文件，爬虫的go文件留一个XXrun()的函数给外部调用，XXrun()会登记在crawlsList.go里面的crawlListStruct结构体中

运行时，用for循环对crawlListStruct的funcList进行遍历，用goroutine对每个XXrun()进行运行

数据库保存一天的数据，每天晚上10点删除数据库原来的数据，然后开始写入新的数据，删除数据不在这个目录内处理，有一个外部统筹crawl和quote的模块的程序

#### 爬取的站点

1，[新华社](http://www.xinhuanet.com/)         ID : 1

2，[中国网](http://www.china.com.cn/)          ID : 2

3，[新浪财经](https://finance.sina.com.cn/)    ID : 3

4，[国家统计局](http://wap.stats.gov.cn/jd/201912/t20191210_1716707.html) ID : 4

5，[第一财经](https://www.yicai.com/)          ID : 5

6，[中国人民共和国中央人民政府](http://www.gov.cn/index.htm)                ID : 6

7，[cnBeta.com](https://www.cnbeta.com/)      ID : 7

8, [鲸鱼钱包](https://whale-alert.io/)         ID: 8
