interface PaginationProps {
  currentPage: number
  totalPages: number
  totalCount: number
  noun: string
  onChange: (page: number) => void
}

export default function Pagination({ currentPage, totalPages, totalCount, noun, onChange }: PaginationProps) {
  if (totalPages <= 1) return null

  return (
    <div className="mt-6 flex items-center justify-between">
      <button
        type="button"
        onClick={() => onChange(Math.max(1, currentPage - 1))}
        disabled={currentPage === 1}
        className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40 transition-colors"
      >
        ← Previous
      </button>
      <span className="text-xs text-gray-500 tabular-nums">
        {currentPage} / {totalPages} &nbsp;·&nbsp; {totalCount} {noun}
      </span>
      <button
        type="button"
        onClick={() => onChange(Math.min(totalPages, currentPage + 1))}
        disabled={currentPage === totalPages}
        className="rounded-lg border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-40 transition-colors"
      >
        Next →
      </button>
    </div>
  )
}
