import { useCallback, useEffect, useState } from 'react'
import ReactFlow, {
  addEdge,
  Background,
  Connection,
  Controls,
  Edge,
  Node,
  OnConnect,
  OnEdgesChange,
  OnNodesChange,
  useEdgesState,
  useNodesState,
  MarkerType,
} from 'reactflow'
import 'reactflow/dist/style.css'
import { v4 as uuid } from 'uuid'
import { fetchTopology, updateTopology, fetchTopologyStatus, fetchPodLogs } from './api'
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

type PodInfo = { podName: string; phase: string; ready: boolean }

export function TopologyCanvas({ topologyId, onSaved }: TopologyCanvasProps) {
  const [nodes, setNodes, onNodesChange] = useNodesState([])
  const [edges, setEdges, onEdgesChange] = useEdgesState([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [status, setStatus] = useState<{ pods: PodInfo[] } | null>(null)
  const [statusLoading, setStatusLoading] = useState(false)
  const [logsModal, setLogsModal] = useState<{ podName: string; container: string } | null>(null)
  const [logs, setLogs] = useState('')
  const [logsLoading, setLogsLoading] = useState(false)

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

  const loadStatus = useCallback(async () => {
    setStatusLoading(true)
    try {
      const s = await fetchTopologyStatus(topologyId)
      setStatus(s?.pods?.length ? { pods: s.pods } : null)
    } catch {
      setStatus(null)
    } finally {
      setStatusLoading(false)
    }
  }, [topologyId])

  const openLogs = useCallback(async (podName: string, container: string) => {
    setLogsModal({ podName, container })
    setLogsLoading(true)
    setLogs('')
    try {
      const text = await fetchPodLogs(topologyId, podName, container)
      setLogs(text)
    } catch (e) {
      setLogs(String(e))
    } finally {
      setLogsLoading(false)
    }
  }, [topologyId])

  const closeLogs = useCallback(() => setLogsModal(null), [])

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
        <button onClick={loadStatus} disabled={statusLoading} title="Обновить статус и список подов">
          {statusLoading ? '...' : 'Статус узлов'}
        </button>
      </div>

      {status && status.pods.length > 0 && (
        <div className="pods-panel">
          <h4>Поды (узлы)</h4>
          <ul className="pods-list">
            {status.pods.map((p) => (
              <li key={p.podName} className="pod-item">
                <span className="pod-name">{p.podName}</span>
                <span className={`pod-phase ${p.phase.toLowerCase()}`}>{p.phase}</span>
                {p.ready && <span className="pod-ready">Ready</span>}
                <span className="pod-logs">
                  <button type="button" onClick={() => openLogs(p.podName, 'ipfs')}>Логи IPFS</button>
                  <button type="button" onClick={() => openLogs(p.podName, 'ipfs-cluster')}>Логи Cluster</button>
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {logsModal && (
        <div className="logs-modal-overlay" onClick={closeLogs} role="presentation">
          <div className="logs-modal" onClick={(e) => e.stopPropagation()} role="dialog">
            <div className="logs-modal-header">
              <h3>Логи {logsModal.podName} ({logsModal.container})</h3>
              <button type="button" onClick={closeLogs} className="logs-close">×</button>
            </div>
            <div className="logs-modal-body">
              {logsLoading ? (
                <p>Загрузка...</p>
              ) : (
                <pre className="logs-content">{logs || '(пусто)'}</pre>
              )}
            </div>
          </div>
        </div>
      )}

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
