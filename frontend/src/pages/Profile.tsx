import React, { useState } from 'react'
import Layout from '../components/Layout'
import { useAuth } from '../auth/AuthContext'
import { apiClient } from '../api/client'
import { getApiErrorMessage } from '../api/errors'

interface FieldSectionProps {
  title: string
  onSave: (value: string) => Promise<void>
  inputType?: string
  placeholder?: string
  currentValue?: string
  confirmLabel?: string
}

function FieldSection({ title, onSave, inputType = 'text', placeholder, currentValue, confirmLabel }: FieldSectionProps) {
  const [value, setValue] = useState('')
  const [confirm, setConfirm] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)

    if (confirmLabel && value !== confirm) {
      setError('Values do not match')
      return
    }

    if (!value.trim()) return

    setLoading(true)
    try {
      await onSave(value)
      setSuccess(true)
      setValue('')
      setConfirm('')
    } catch (err) {
      setError(getApiErrorMessage(err, 'Failed to update'))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-6">
      <h2 className="text-lg font-semibold text-gray-900 mb-1">{title}</h2>
      {currentValue && (
        <p className="text-sm text-gray-500 mb-4">Current: <span className="font-medium text-gray-700">{currentValue}</span></p>
      )}
      <form onSubmit={handleSubmit} className="space-y-3">
        <input
          type={inputType}
          value={value}
          onChange={e => { setValue(e.target.value); setSuccess(false); setError('') }}
          placeholder={placeholder || `New ${title.toLowerCase()}`}
          className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        {confirmLabel && (
          <input
            type={inputType}
            value={confirm}
            onChange={e => { setConfirm(e.target.value); setError('') }}
            placeholder={confirmLabel}
            className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        )}
        {error && <p className="text-sm text-red-600">{error}</p>}
        {success && <p className="text-sm text-green-600">Updated successfully</p>}
        <button
          type="submit"
          disabled={loading || !value.trim()}
          className="px-4 py-2 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? 'Saving...' : 'Save changes'}
        </button>
      </form>
    </div>
  )
}

export default function Profile() {
  const { user, refresh } = useAuth()

  const updateField = async (patch: Record<string, string>) => {
    await apiClient.put('/auth/me', patch)
    await refresh()
  }

  return (
    <Layout>
      <div className="max-w-lg mx-auto py-10 px-4">
        <h1 className="text-2xl font-bold text-gray-900 mb-8">Profile</h1>
        <div className="space-y-4">
          <FieldSection
            title="Username"
            currentValue={user?.username}
            placeholder="New username"
            onSave={value => updateField({ username: value })}
          />
          <FieldSection
            title="Email"
            inputType="email"
            currentValue={user?.email}
            placeholder="New email address"
            onSave={value => updateField({ email: value })}
          />
          <FieldSection
            title="Password"
            inputType="password"
            placeholder="New password"
            confirmLabel="Confirm new password"
            onSave={value => updateField({ new_password: value })}
          />
        </div>
      </div>
    </Layout>
  )
}
