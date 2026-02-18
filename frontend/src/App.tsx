import { useState, useEffect } from 'react'
import { fetchTopologies, createTopology, deleteTopology, deployTopology, undeployTopology, fetchTopologyStatus } from './api'
import type { TopologySummary } from './types'
import { TopologyCanvas } from './TopologyCanvas'
import './App.css'

function App() {
  const [topologies, setTopologies] = useState<TopologySummary[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedId, setSelectedId] = useState<string | null>(null)
  const [creating, setCreating] = useState(false)
  const [newName, setNewName] = useState('')

  const loadTopologies = async () => {
    try {
      const list = await fetchTopologies()
      setTopologies(list)
    } catch (e) {
      console.error(e)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    loadTopologies()
  }, [])

  const handleCreate = async () => {
    if (!newName.trim()) return
    setCreating(true)
    try {
      const t = await createTopology({ name: newName.trim(), nodes: [], edges: [] })
      setNewName('')
      await loadTopologies()
      setSelectedId(t.topologyId)
    } catch (e) {
      console.error(e)
    } finally {
      setCreating(false)
    }
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteTopology(id)
      if (selectedId === id) setSelectedId(null)
      await loadTopologies()
    } catch (e) {
      console.error(e)
    }
  }

  const handleDeploy = async (id: string) => {
    try {
      await deployTopology(id, 'default')
      await loadTopologies()
    } catch (e) {
      console.error(e)
    }
  }

  const handleUndeploy = async (id: string) => {
    try {
      await undeployTopology(id)
      await loadTopologies()
    } catch (e) {
      console.error(e)
    }
  }

  const selected = topologies.find((t) => t.topologyId === selectedId)

  return (
    <div className="app">
      <header className="app-header">
        <h1>IPFS Cluster Topology Orchestrator</h1>
        <p>Создайте топологию сети: добавьте узлы и соедините их. Ребро A → B означает: узел A бустрапится к B (B — bootstrap).</p>
      </header>

      <div className="app-body">
        <aside className="sidebar">
          <section className="create-section">
            <input
              type="text"
              placeholder="Название топологии"
              value={newName}
              onChange={(e) => setNewName(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && handleCreate()}
            />
            <button onClick={handleCreate} disabled={creating || !newName.trim()}>
              {creating ? '...' : 'Создать'}
            </button>
          </section>

          <section className="list-section">
            <h3>Топологии</h3>
            {loading ? (
              <p>Загрузка...</p>
            ) : topologies.length === 0 ? (
              <p className="empty">Нет топологий</p>
            ) : (
              <ul>
                {topologies.map((t) => (
                  <li key={t.topologyId} className={selectedId === t.topologyId ? 'selected' : ''}>
                    <button className="item" onClick={() => setSelectedId(t.topologyId)}>
                      <span className="name">{t.name}</span>
                      <span className="meta">{t.nodeCount} нод · {t.edgeCount} рёбер</span>
                      <span className={`status status-${t.deployStatus}`}>{t.deployStatus}</span>
                    </button>
                    <div className="actions">
                      {t.deployStatus === 'none' || t.deployStatus === 'error' ? (
                        <button onClick={() => handleDeploy(t.topologyId)} title="Деплой">🚀</button>
                      ) : (
                        <button onClick={() => handleUndeploy(t.topologyId)} title="Undeploy">⏹</button>
                      )}
                      <button onClick={() => handleDelete(t.topologyId)} title="Удалить">🗑</button>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </section>
        </aside>

        <main className="canvas-area">
          {selected ? (
            <TopologyCanvas
              topologyId={selected.topologyId}
              onSaved={loadTopologies}
            />
          ) : (
            <div className="placeholder">
              <p>Выберите топологию слева или создайте новую</p>
            </div>
          )}
        </main>
      </div>
    </div>
  )
}

export default App
