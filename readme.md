
#本文档已过时，以代码为准。

#设计要点  
1. 合理性：只有一个坏球，这个球的重量不变。  
1. 唯一性：同构的方案（仅仅改变编号就成为相同方案）只保留一个。  
1. 有效性：能有效确定坏球以及它是偏轻还是偏重。

1. 方案的各步骤不是唯一的，一个称量结果可以有多套合格的后续步骤。因此递归时不能简单地找到一种方案就回溯。


#掩码  
1. 称量结果：包括一次称量摆位和称量结果。  
1. 可能集：12个球的可能性集合{标准球、偏重球、偏轻球}。  
1. 掩码：每个称量结果缩小可能集，不能出现的可能性就是这个称量结果的掩码。  

#数据格式
1. 可能集：2个整数分别表示偏重、偏轻的可能性。各使用最低12位，每1位是一个球的状态。
1. 掩码：2个整数，位数对应于可能集一样。1表示未消除的可能性，0表示称量结果所不兼容的可能。  
1. 可能集与掩码做与运算，结果是新的可能集。三次称量之后只剩下一个bit是1，则方案成功。 
 


