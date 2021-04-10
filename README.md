# go-truffle
web frame core for go gin

# consul refer
https://github.com/zhangyunan1994/lemon

# feign refer
https://github.com/HikoQiu/go-feign


# go接口赋值原则
1. 看右边是不是左边的子类， 是子类即可赋值
2. 如何判断是否子类：看这个类是否实现了接口的所有方法
3. T 与 *T 不是一个类型。 类型*T的方法集包含所有 receiver T + *T方法，类型T的方法集只包含所有 receiver T方法
