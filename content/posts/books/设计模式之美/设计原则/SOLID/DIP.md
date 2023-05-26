---
title: "DIP"
date: 2023-05-23T18:30:01+08:00
summary: Dependency Inversion Principle 依赖反转原则/依赖倒置原则
---

# IOC

Inversion Of Control, 控制反转

- 流程的控制权从程序员“反转”到了框架。
- 控制反转并不是一种具体的实现技巧，而是一个比较笼统的设计思想，一般用来指导框架层面的设计。

# DI

Dependency Injection, 依赖注入

- 一种具体的编码技巧
- 不通过 new() 的方式在类内部创建依赖类对象，而是将依赖的类对象在外部创建好之后，通过构造函数、函数参数等方式传递（或注入）给类使用。

优点

- 提高了代码的扩展性，我们可以灵活地替换依赖的类。

# DIP

Dependency Inversion Principle, 依赖反转原则/依赖倒置原则

高层模块（high-level modules）不要依赖低层模块（low-level）。高层模块和低层模块应该通过抽象（abstractions）来互相依赖。除此之外，抽象
（abstractions）不要依赖具体实现细节（details），具体实现细节（details）依赖抽象（abstractions）。

- 调用者属于高层，被调用者属于低层
