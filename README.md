## 最近看了下压缩算法 自己写了个基于huffman的压缩算法(demo)，这只是zip压缩算法的一部分

## 参考文章
> - [ZIP压缩算法详细分析及解压实例解释](https://www.cnblogs.com/esingchan/p/3958962.html)
> - [详细图解哈夫曼Huffman编码树](https://blog.csdn.net/FX677588/article/details/70767446)
> - [几种压缩算法实现原理详解](https://blog.csdn.net/ghevinn/article/details/45747465)

### 2018-06-11  demo版

> - 把算法基本完成了，但是压缩了不是很高，简单看下问题，是因为序列化huffman树的占用了很大的空间 这里有很大的优化空间（目前是用byte数组存放的 可以把这个数组用比特位表示）
> - 对于中文的压缩率没有英文高  因为这个是基于字节去建立huffman树的 中文编码中的byte分母要比英文大很多
