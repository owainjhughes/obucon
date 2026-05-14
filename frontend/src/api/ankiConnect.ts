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
    const host = window.location.hostname
    const onLocalhost = host === "localhost" || host === "127.0.0.1"
    if (onLocalhost) {
      throw new Error("Could not reach Anki. Make sure Anki is open and the AnkiConnect plugin is installed.")
    }
    throw new Error(
      `Could not reach Anki from ${window.location.origin}. Either Anki is not running, or AnkiConnect has not whitelisted this site. ` +
        `In Anki: Tools → Add-ons → AnkiConnect → Config, and add "${window.location.origin}" to "webCorsOriginList", then restart Anki.`,
    )
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

function stripHtml(value: string): string {
  return value.replace(/<[^>]*>/g, "").trim()
}

export interface AnkiDeckPreview {
  fields: string[]
  sample: Record<string, string>
  noteCount: number
}

export async function getAnkiDeckNames(): Promise<string[]> {
  return invoke<string[]>("deckNames")
}

const MATURE_QUERY = (deckName: string) => `deck:"${deckName}" "prop:ivl>=21"`

export async function previewAnkiDeck(deckName: string): Promise<AnkiDeckPreview> {
  const noteIds = await invoke<number[]>("findNotes", { query: MATURE_QUERY(deckName) })
  if (noteIds.length === 0) {
    return { fields: [], sample: {}, noteCount: 0 }
  }

  const notes = await invoke<AnkiNoteInfo[]>("notesInfo", { notes: [noteIds[0]] })
  const first = notes[0]
  if (!first) {
    return { fields: [], sample: {}, noteCount: noteIds.length }
  }

  const sortedEntries = Object.entries(first.fields).sort((a, b) => a[1].order - b[1].order)
  const fields = sortedEntries.map(([name]) => name)
  const sample: Record<string, string> = {}
  for (const [name, info] of sortedEntries) {
    sample[name] = stripHtml(info.value)
  }

  return { fields, sample, noteCount: noteIds.length }
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
  lemmaField: string,
  meaningField: string,
  hiraganaField?: string,
): Promise<Array<{ lemma: string; meaning: string; hiragana: string }>> {
  const noteIds = await invoke<number[]>("findNotes", { query: MATURE_QUERY(deckName) })
  if (noteIds.length === 0) {
    return []
  }

  const notes = await invoke<AnkiNoteInfo[]>("notesInfo", { notes: noteIds })

  return notes
    .map((note) => {
      const lemmaRaw = note.fields[lemmaField]?.value ?? ""
      const meaningRaw = note.fields[meaningField]?.value ?? ""
      const hiraganaRaw = hiraganaField ? note.fields[hiraganaField]?.value ?? "" : ""
      return {
        lemma: stripHtml(lemmaRaw),
        meaning: stripHtml(meaningRaw),
        hiragana: stripHtml(hiraganaRaw),
      }
    })
    .filter((e) => e.lemma !== "")
}
