import { useCallback, useEffect, useState } from 'react'
import ReactFlow, {
  addEdge,
  Background,
  Connection,
  Controls,
  Edge,
  Node,
  NodeChange,
  OnConnect,
  OnEdgesChange,
  OnNodesChange,
  useEdgesState,
  useNodesState,
  MarkerType,
} from 'reactflow'
import 'reactflow/dist/style.css'
import { v4 as uuid } from 'uuid'
import { fetchTopology, updateTopology } from './api'
import type { TopologyNode, TopologyEdge } from './types'

const nodeTypes = {} as const

function toFlowNode(n: TopologyNode): Node {
  return {
    id: n.nodeId,
    type: 'default',
    position: n.position,
    data: { label: n.label },
  }
}

function toFlowEdge(e: TopologyEdge): Edge {
  return {
    id: e.edgeId,
    source: e.sourceNodeId,
    target: e.targetNodeId,
    markerEnd: { type: MarkerType.ArrowClosed },
  }
}

interface TopologyCanvasProps {
  topologyId: string
  onSaved?: () => void
}

export function TopologyCanvas({ topologyId, onSaved }: TopologyCanvasProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState([])
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const t = await fetchTopology(topologyId)
      setNodes(t.nodes.map(toFlowNode))
      setEdges(t.edges.map(toFlowEdge))
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }, [topologyId, setNodes, setEdges])

  useEffect(() => {
    load()
  }, [load])

  const onConnect: OnConnect = useCallback(
    (conn: Connection) => {
      if (!conn.source || !conn.target) return
      setEdges((eds) =>
        addEdge({ ...conn, id: uuid() } as Edge, eds)
      )
    },
    [setEdges]
  )

  const handleSave = useCallback(async () => {
    const topologyNodes: TopologyNode[] = nodes.map((n) => ({
      nodeId: n.id,
      label: (n.data?.label as string) || n.id,
      position: n.position,
      role: 'worker' as const,
    }))
    const topologyEdges: TopologyEdge[] = edges.map((e) => ({
      edgeId: e.id,
      sourceNodeId: e.source!,
      targetNodeId: e.target!,
    }))
    setSaving(true)
    try {
      await updateTopology(topologyId, { nodes: topologyNodes, edges: topologyEdges })
      onSaved?.()
    } catch (e) {
      console.error(e)
    } finally {
      setSaving(false)
    }
  }, [topologyId, nodes, edges, onSaved])

  const handleAddNode = useCallback(() => {
    const id = uuid()
    setNodes((nds) => [
      ...nds,
      {
        id,
        type: 'default',
        position: { x: 250 + Math.random() * 100, y: 150 + Math.random() * 100 },
        data: { label: `Node ${nds.length + 1}` },
      },
    ])
  }, [setNodes])

  if (loading) {
    return (
      <div className="canvas-loading">
        <p>Загрузка...</p>
      </div>
    )
  }

  return (
    <div className="canvas-wrapper">
      <div className="canvas-toolbar">
        <button onClick={handleAddNode}>+ Добавить узел</button>
        <button onClick={handleSave} disabled={saving}>
          {saving ? 'Сохранение...' : 'Сохранить'}
        </button>
      </div>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange as OnNodesChange}
        onEdgesChange={onEdgesChange as OnEdgesChange}
        onConnect={onConnect}
        fitView
        nodeTypes={nodeTypes}
      >
        <Background color="#30363d" gap={16} />
        <Controls />
      </ReactFlow>
    </div>
  )
}
