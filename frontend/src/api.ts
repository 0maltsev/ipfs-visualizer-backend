const API = '/v1'

export async function fetchTopologies(): Promise<any[]> {
  const r = await fetch(`${API}/topologies`)
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}

export async function fetchTopology(id: string): Promise<any> {
  const r = await fetch(`${API}/topologies/${id}`)
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}

export async function createTopology(data: { name: string; nodes?: any[]; edges?: any[] }): Promise<any> {
  const r = await fetch(`${API}/topologies`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}

export async function updateTopology(id: string, data: { name?: string; nodes?: any[]; edges?: any[] }): Promise<any> {
  const r = await fetch(`${API}/topologies/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  })
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}

export async function deleteTopology(id: string): Promise<void> {
  const r = await fetch(`${API}/topologies/${id}`, { method: 'DELETE' })
  if (!r.ok) throw new Error(await r.text())
}

export async function deployTopology(id: string, namespace?: string): Promise<any> {
  const r = await fetch(`${API}/topologies/${id}/deploy`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ namespace: namespace || 'default' }),
  })
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}

export async function undeployTopology(id: string): Promise<void> {
  const r = await fetch(`${API}/topologies/${id}/undeploy`, { method: 'POST' })
  if (!r.ok) throw new Error(await r.text())
}

export async function fetchTopologyStatus(id: string): Promise<any> {
  const r = await fetch(`${API}/topologies/${id}/status`)
  if (!r.ok) throw new Error(await r.text())
  return r.json()
}
