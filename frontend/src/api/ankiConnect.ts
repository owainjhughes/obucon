const ANKI_CONNECT_URL = "http://127.0.0.1:8765" // default anki connect URL. Could make this user defined later on?
const ANKI_CONNECT_VERSION = 6

interface AnkiResponse<T> {
  result: T
  error: string | null
}

async function invoke<T>(action: string, params?: Record<string, unknown>): Promise<T> {
  const body: Record<string, unknown> = { action, version: ANKI_CONNECT_VERSION }
  if (params !== undefined) {
    body.params = params
  }

  let response: Response
  try {
    response = await fetch(ANKI_CONNECT_URL, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    })
  } catch {
    throw new Error("Could not reach Anki. Make sure Anki is open and the AnkiConnect plugin is installed.")
  }

  if (!response.ok) {
    throw new Error(`AnkiConnect returned HTTP ${response.status}`)
  }

  const data: AnkiResponse<T> = await response.json()

  // if error type is string, there is one error. if there is many errors, it is an error array. error array is scary!
  if (typeof data.error === "string") {
    throw new Error(data.error)
  }
  return data.result
}

interface AnkiNoteInfo {
  noteId: number
  fields: Record<string, { value: string; order: number }>
  tags: string[]
}

export async function getAnkiDeckNames(): Promise<string[]> {
  return invoke<string[]>("deckNames")
}

export async function exportVocabToAnki(
  entries: Array<{ lemma: string; meaning: string }>,
  deckName: string,
): Promise<{ added: number; skipped: number }> {
  await invoke<number>("createDeck", { deck: deckName })

  const notes = entries.map((entry) => ({
    deckName,
    modelName: "Basic",
    fields: {
      Front: entry.lemma,
      Back: entry.meaning || "",
    },
    options: { allowDuplicate: true },
    tags: ["ginapi-vocab"],
  }))

  const results = await invoke<Array<number | null>>("addNotes", { notes })
  const added = results.filter((id) => id !== null).length
  const skipped = results.filter((id) => id === null).length
  return { added, skipped }
}

export async function importVocabFromAnki(
  deckName: string,
): Promise<Array<{ lemma: string; meaning: string }>> {
  const noteIds = await invoke<number[]>("findNotes", { query: `deck:"${deckName}"` })
  if (noteIds.length === 0) {
    return []
  }

  const notes = await invoke<AnkiNoteInfo[]>("notesInfo", { notes: noteIds })

  return notes
    .map((note) => {
      const sorted = Object.entries(note.fields).sort((a, b) => a[1].order - b[1].order)
      const lemma = (sorted[0]?.[1].value ?? "").replace(/<[^>]*>/g, "").trim()
      const meaning = (sorted[1]?.[1].value ?? "").replace(/<[^>]*>/g, "").trim()
      return { lemma, meaning }
    })
    .filter((e) => e.lemma !== "")
}
