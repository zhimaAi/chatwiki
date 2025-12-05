import ELK from 'elkjs'

const elk = new ELK()

// 定义一个简单的图形结构
const graph = {
  id: 'root',
  children: [
    {
      id: 'parent',
      children: [
        { id: 'c1', width: 40, height: 40 },
        { id: 'c2', width: 40, height: 40 }
      ],
      width: 200,
      height: 150
    },
    { id: 'n3', width: 50, height: 50 }
  ],
  edges: [{ id: 'e1', sources: ['c1'], targets: ['n3'] }]
}

// 计算布局
elk
  .layout(graph)
  .then((layoutedGraph) => {
    console.log(layoutedGraph)
  })
  .catch(console.error)