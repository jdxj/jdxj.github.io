---
title: "LSP"
date: 2023-05-24T11:35:44+08:00
weight: 3
summary: Liskov Substitution Principle 里式替换原则
---

Liskov Substitution Principle, 里式替换原则

子类对象（object of subtype/derived class）能够替换程序（program）中父类对象（object of base/parent class）出现的任何地方，并且保证
原来程序的逻辑行为（behavior）不变及正确性不被破坏。

LSP与多态的区别

- 里式替换是一种设计原则，是用来指导继承关系中子类该如何设计的，子类的设计要保证在替换父类的时候，不改变原有程序的逻辑以及不破坏原有程序的正确性。

> 替换前后的行为一模一样
