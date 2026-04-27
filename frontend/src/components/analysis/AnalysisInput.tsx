import React, { useState } from 'react'
import { apiClient } from '../../api/client'
import { getApiErrorMessage } from '../../api/errors'
import AnalysisOutput from './AnalysisOutput'

type AnalysisMode = 'text' | 'file'

interface Token {
  surface: string
  lemma: string
  pos: string
  is_known: boolean
  grade_level?: number | null
}

interface AnalysisResult {
  tokens: Token[]
  total_tokens: number
  missing: string[]
}

export default function AnalysisInput() {
  const [mode, setMode] = useState<AnalysisMode>('text')
  const [text, setText] = useState('')
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [language] = useState('ja')
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState('')
  const [result, setResult] = useState<AnalysisResult | null>(null)

  const handleAnalyse = async () => {
    setError('')
    setResult(null)
    setIsLoading(true)

    try {
      let response

      if (mode === 'file') {
        if (!selectedFile) {
          setError('Please choose a .txt, .md, .docx, or .pdf file first')
          setIsLoading(false)
          return
        }

        const payload = new FormData()
        payload.append('file', selectedFile)
        payload.append('language', language)

        response = await apiClient.post('/analyze/file', payload, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        })
      } else {
        response = await apiClient.post('/analyze', {
          text,
          language,
        })
      }

      setResult(response.data)
    } catch (err: unknown) {
      setError(getApiErrorMessage(err, 'Analysis failed'))
    } finally {
      setIsLoading(false)
    }
  }

  const resetAnalysis = () => {
    setResult(null)
    setError('')
  }

  const canAnalyse = mode === 'file' ? Boolean(selectedFile) : Boolean(text.trim())

  if (result) {
    return (
      <AnalysisOutput
        tokens={result.tokens}
        missing={result.missing}
        onReset={resetAnalysis}
      />
    )
  }

  return (
    <section className="px-4 py-10">
      <div className="mx-auto rounded-xl border border-gray-200 bg-white p-6 shadow-sm max-w-3xl">
        <h2 className="text-lg font-semibold text-gray-900">Text Analysis</h2>
        <p className="mt-1 text-sm text-gray-500">Choose a mode to provide your input.</p>

        <div className="mt-4 inline-flex rounded-lg border border-gray-200 bg-gray-50 p-1 gap-1">
          <button
            type="button"
            onClick={() => setMode('text')}
            className={`rounded-md px-4 py-1.5 text-sm font-medium transition-colors ${
              mode === 'text'
                ? 'bg-white text-[#55F] shadow-sm border border-gray-200'
                : 'text-gray-500 hover:text-gray-900'
            }`}
          >
            Text
          </button>
          <button
            type="button"
            onClick={() => setMode('file')}
            className={`rounded-md px-4 py-1.5 text-sm font-medium transition-colors ${
              mode === 'file'
                ? 'bg-white text-[#55F] shadow-sm border border-gray-200'
                : 'text-gray-500 hover:text-gray-900'
            }`}
          >
            File
          </button>
        </div>

        <div className="mt-6">
          {mode === 'text' && (
            <label className="block">
              <span className="text-sm font-medium text-gray-600">Paste text</span>
              <textarea
                rows={8}
                value={text}
                onChange={(e) => setText(e.target.value)}
                className="mt-2 w-full rounded-lg border border-gray-300 p-3 text-sm text-gray-900 focus:border-[#55F] focus:outline-none focus:ring-2 focus:ring-[#55F]/10 resize-none"
                placeholder="Paste the text you want to analyse"
              />
            </label>
          )}

          {mode === 'file' && (
            <div className="rounded-xl border-2 border-dashed border-gray-200 bg-gray-50 px-6 py-8 text-center hover:border-[#55F] hover:bg-indigo-50/30 transition-colors">
              <label className="block cursor-pointer">
                <span className="text-sm font-medium text-gray-600">Upload file</span>
                <span className="block text-xs text-gray-400 mt-0.5">.txt, .md, .docx, .pdf</span>
                <input
                  type="file"
                  accept=".txt,.md,.docx,.pdf,text/plain,text/markdown,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/pdf"
                  onChange={(e) => {
                    const file = e.target.files && e.target.files.length > 0 ? e.target.files[0] : null
                    setSelectedFile(file)
                    setError('')
                  }}
                  className="mt-3 block w-full text-sm text-gray-500 file:mr-3 file:cursor-pointer file:rounded-lg file:border-0 file:bg-[#55F] file:px-4 file:py-2 file:text-sm file:font-medium file:text-white hover:file:bg-[#44E]"
                />
              </label>
              {selectedFile && (
                <p className="mt-3 text-sm text-gray-600">
                  Selected: <span className="font-medium text-gray-800">{selectedFile.name}</span> ({Math.max(1, Math.round(selectedFile.size / 1024))} KB)
                </p>
              )}
            </div>
          )}

          <div className="mt-4 flex items-center justify-end gap-3">
            <button
              onClick={handleAnalyse}
              disabled={isLoading || !canAnalyse}
              className="rounded-lg border border-[#55F] bg-[#55F] px-5 py-2 text-sm font-semibold text-white hover:bg-[#44E] disabled:cursor-not-allowed disabled:opacity-50 transition-colors"
            >
              {isLoading ? 'Analysing…' : 'Analyse'}
            </button>
          </div>

          {error && (
            <div className="mt-4 rounded-lg border border-red-200 bg-red-50 px-4 py-3 text-sm text-red-700">
              {error}
            </div>
          )}
        </div>
      </div>
    </section>
  )
}
