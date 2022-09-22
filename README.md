# goconvert

## ConvertNilSlice2Empty

### 使用场景

前端组件不支持数据为null的slice，因此需要将为null的slice转换为[]。

如果只是简单的数据结构，那么直接使用make初始化切片即可，但是在复杂的数据结构中，切片数据往往嵌套很深，这时候对每个切片类型进行make操作会很麻烦，因此封装了该函数来实现上述需求。
