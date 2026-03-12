import React, { useEffect, useState } from "react"
import Layout from "../components/Layout"
import { apiClient } from "../api/client"
import { getApiErrorMessage } from "../api/errors"

interface VocabEntry {
  lemma: string
  grade_level?: number | null
  meaning: string
}

export default function Vocab() {
  const [vocab, setVocab] = useState<VocabEntry[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState("")
  const [jlptLevel, setJlptLevel] = useState("N5")
  const [importing, setImporting] = useState(false)

  const loadVocab = async () => {
    setIsLoading(true)
    setError(null)

    try {
      const response = await apiClient.get("/vocab")
      setVocab(response.data.vocab || [])
    } catch (err: unknown) {
      setError(getApiErrorMessage(err, "Failed to load vocabulary"))
    } finally {
      setIsLoading(false)
    }
  }

  useEffect(() => {
    loadVocab()
  }, [])

  const filteredVocab = vocab.filter((entry) =>
    entry.lemma.includes(search) || entry.meaning.toLowerCase().includes(search.toLowerCase())
  )

  const handleImport = async () => {
    setImporting(true)
    setError(null)

    try {
      await apiClient.post("/vocab/bulk", {
        jlpt_level: jlptLevel,
        language: "ja",
      })
      await loadVocab()
    } catch (err: unknown) {
      setError(getApiErrorMessage(err, "Failed to import vocabulary"))
    } finally {
      setImporting(false)
    }
  }

  return (
    <Layout>
      <section className="mx-auto max-w-4xl px-4 py-10">
        <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
          <h1 className="text-xl font-semibold text-gray-900">Known Vocabulary</h1>
          <p className="mt-1 text-sm text-gray-600">
            Words you have marked as known. Use this list to review or update your vocabulary.
          </p>

          <div className="mt-6 flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
            <div className="flex flex-1 flex-col gap-2 md:flex-row md:items-center">
              <input
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                placeholder="Search by word or meaning"
                className="w-full rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm text-gray-900 focus:border-[#55F] focus:outline-none md:max-w-sm"
              />
              <div className="flex items-center gap-2">
                <select
                  value={jlptLevel}
                  onChange={(e) => setJlptLevel(e.target.value)}
                  className="rounded-lg border border-gray-300 bg-white px-3 py-2 text-sm text-gray-900 focus:border-[#55F] focus:outline-none"
                >
                  <option value="N5">JLPT N5</option>
                  <option value="N4">JLPT N4</option>
                  <option value="N3">JLPT N3</option>
                  <option value="N2">JLPT N2</option>
                  <option value="N1">JLPT N1</option>
                </select>
                <button
                  type="button"
                  onClick={handleImport}
                  disabled={importing}
                  className="rounded-full border border-[#55F] bg-[#55F] px-4 py-2 text-sm font-semibold text-white hover:bg-[#44E] disabled:cursor-not-allowed disabled:opacity-50"
                >
                  {importing ? "Importing..." : "Import JLPT list"}
                </button>
              </div>
            </div>
          </div>

          {isLoading ? (
            <div className="mt-8 text-sm text-gray-600">Loading...</div>
          ) : error ? (
            <div className="mt-8 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
              {error}
            </div>
          ) : (
            <div className="mt-6 overflow-x-auto">
              <table className="w-full min-w-[520px] table-auto border-collapse">
                <thead>
                  <tr className="bg-gray-50 text-left text-xs font-medium uppercase tracking-wider text-gray-600">
                    <th className="px-4 py-3">Word</th>
                    <th className="px-4 py-3">Meaning</th>
                    <th className="px-4 py-3">Grade</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-100 text-sm text-gray-700">
                  {filteredVocab.map((entry) => (
                    <tr key={entry.lemma}>
                      <td className="px-4 py-3 font-semibold text-gray-900">{entry.lemma}</td>
                      <td className="px-4 py-3">{entry.meaning || "—"}</td>
                      <td className="px-4 py-3">{entry.grade_level ?? "—"}</td>
                    </tr>
                  ))}
                </tbody>
              </table>

              {filteredVocab.length === 0 && (
                <div className="mt-6 rounded-lg border border-gray-200 bg-gray-50 px-4 py-3 text-sm text-gray-600">
                  No known words found.
                </div>
              )}
            </div>
          )}
        </div>
      </section>
    </Layout>
  )
}
