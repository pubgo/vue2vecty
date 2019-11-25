## 文本
<span>Message: {{ msg }}</span>
<span v-once>这个将不会改变: {{ msg }}</span>

## 原始 HTML
<p>Using mustaches: {{ rawHtml }}</p>
<p>Using v-html directive: <span v-html="rawHtml"></span></p>

## 特性
<div v-bind:id="dynamicId"></div>
<button v-bind:disabled="isButtonDisabled">Button</button>
如果 isButtonDisabled 的值是 null、undefined 或 false，则 disabled 特性甚至不会被包含在渲染出来的 <button> 元素中。

## 使用 JavaScript 表达式
{{ number + 1 }}
{{ ok ? 'YES' : 'NO' }}
{{ message.split('').reverse().join('') }}
<div v-bind:id="'list-' + id"></div>

## 指令
<p v-if="seen">现在你看到我了</p>

## 参数
<a v-bind:href="url">...</a>
<a v-on:click="doSomething">...</a>

## v-bind 缩写
<!-- 完整语法 -->
<a v-bind:href="url">...</a>

<!-- 缩写 -->
<a :href="url">...</a>

## v-on 缩写
<!-- 完整语法 -->
<a v-on:click="doSomething">...</a>

<!-- 缩写 -->
<a @click="doSomething">...</a>

## 计算属性
<div id="example">
  {{ message.split('').reverse().join('') }}
</div>

## 计算属性 vs 侦听属性
var vm = new Vue({
  el: '#demo',
  data: {
    firstName: 'Foo',
    lastName: 'Bar',
    fullName: 'Foo Bar'
  },
  watch: {
    firstName: function (val) {
      this.fullName = val + ' ' + this.lastName
    },
    lastName: function (val) {
      this.fullName = this.firstName + ' ' + val
    }
  }
})

##  sss
<input v-model="question">
<div v-bind:class="{ active: isActive }"></div>
<div v-bind:class="classObject"></div>
<div v-bind:style="{ color: activeColor, fontSize: fontSize + 'px' }"></div>
<h1 v-if="awesome">Vue is awesome!</h1>
<h1 v-else>Oh no 😢</h1>

<template v-if="ok">
  <h1>Title</h1>
  <p>Paragraph 1</p>
  <p>Paragraph 2</p>
</template>

<div v-if="Math.random() > 0.5">
  Now you see me
</div>
<div v-else>
  Now you don't
</div>