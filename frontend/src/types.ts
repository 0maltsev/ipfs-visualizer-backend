export interface TopologyNode {
  nodeId: string
  label: string
  position: { x: number; y: number }
  role: 'bootstrap' | 'worker'
}

export interface TopologyEdge {
  edgeId: string
  sourceNodeId: string
  targetNodeId: string
}

export interface Topology {
  topologyId: string
  name: string
  nodes: TopologyNode[]
  edges: TopologyEdge[]
  deployStatus: string
  k8sNamespace?: string
  createdAt: string
  updatedAt: string
}

export interface TopologySummary {
  topologyId: string
  name: string
  nodeCount: number
  edgeCount: number
  deployStatus: string
  createdAt: string
}
