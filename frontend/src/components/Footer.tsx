import React from 'react'

export default function Footer() {
    return (
        <footer className="border-t border-gray-200 bg-white">
            <div className="mx-auto max-w-6xl px-6 py-8">
                <div className="flex flex-col items-center text-center">
                    <div className="flex items-center gap-2.5">
                        <img alt="ObuCon" src="/android-chrome-192x192.png" className="h-6 w-auto" />
                        <div className="text-sm font-semibold text-gray-900 tracking-wide">ObuCon</div>
                    </div>
                    <div className="mt-4 flex items-center justify-center gap-4 text-gray-400">
                        <a href="https://ojhdev.pythonanywhere.com/" className="text-xs hover:text-gray-700 transition-colors">Korean Prototype</a>
                        <span className="text-gray-200">|</span>
                        <a href="https://www.linkedin.com/in/ohdev" className="text-xs hover:text-gray-700 transition-colors">LinkedIn</a>
                    </div>
                </div>
                <div className="mt-6 border-t border-gray-100 pt-5 text-center text-xs text-gray-400">
                    © 2026 ObuCon. All rights reserved.
                </div>
            </div>
        </footer>
    )
}
