import React from 'react'

export default function Footer() {
    return (
        <footer className="border-t border-gray-200 bg-white">
            <div className="mx-auto max-w-6xl px-6 py-12">
                    <div className="flex flex-col items-center text-center">
                        <div className="flex items-center gap-3">
							<img alt="ObuCon" src="/android-chrome-192x192.png" className="h-7 w-auto" />
                            <div className="text-lg font-semibold text-gray-900">ObuCon</div>
                        </div>
                        <div className="mt-6 flex items-center justify-center gap-4 text-gray-500">
                            <a href="https://ojhdev.pythonanywhere.com/" className="text-sm hover:text-gray-900">Korean Prototype site</a>
                            <a href="https://www.linkedin.com/in/ohdev" className="text-sm hover:text-gray-900">LinkedIn</a>
                        </div>
                    </div>

                <div className="mt-10 border-t border-gray-200 pt-6 text-center text-xs text-gray-500">
                    © 2026 ObuCon. All rights reserved.
                </div>
            </div>
        </footer>
    )
}
