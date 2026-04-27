import { Disclosure, DisclosureButton, DisclosurePanel, Menu, MenuButton, MenuItem, MenuItems } from "@headlessui/react"
import { Bars3Icon, XMarkIcon } from "@heroicons/react/24/outline"
import { Link, useLocation, useNavigate } from "react-router-dom"
import { useAuth } from "../auth/AuthContext"

const navigation = [
	{ name: "Home", href: "/" },
	{ name: "Analysis", href: "/analysis" },
	{ name: "Vocabulary", href: "/vocabulary" },
	{ name: "Dictionary", href: "/dictionary" },
	{ name: "jWiki", href: "https://ja.wikipedia.org/wiki/%E3%83%A1%E3%82%A4%E3%83%B3%E3%83%9A%E3%83%BC%E3%82%B8", external: true }
]

function classNames(...classes: string[]) {
	return classes.filter(Boolean).join(" ")
}

export default function Navbar() {
	const { user, logout } = useAuth()
	const navigate = useNavigate()
	const location = useLocation()

	const handleLogout = async () => {
		await logout()
		navigate("/login")
	}

	return (
		<Disclosure
			as="nav"
			className="bg-white border-b border-gray-200">
			<div className="mx-auto max-w-6xl px-4">
				<div className="flex h-14 items-stretch gap-2">
					{/* Hamburger - mobile only */}
					<div className="sm:hidden flex items-center">
						<DisclosureButton className="group inline-flex items-center justify-center rounded p-2 text-gray-500 hover:bg-gray-100 hover:text-gray-900">
							<span className="sr-only">Open main menu</span>
							<Bars3Icon aria-hidden="true" className="block size-5 group-data-[open]:hidden" />
							<XMarkIcon aria-hidden="true" className="hidden size-5 group-data-[open]:block" />
						</DisclosureButton>
					</div>
					<div className="flex flex-1 items-stretch">
						<div className="flex shrink-0 items-center">
							<img alt="ObuCon" src="/android-chrome-192x192.png" className="h-7 w-auto" />
						</div>
						<div className="hidden sm:ml-6 sm:flex sm:items-stretch">
							<div className="flex items-stretch gap-1">
								{navigation.map((item) => {
								const isCurrent = !item.external && (item.href === "/"
									? location.pathname === item.href
									: location.pathname.startsWith(item.href))
								if (item.external) {
									return (
										<a
											key={item.name}
											href={item.href}
											target="_blank"
											rel="noopener noreferrer"
											className="inline-flex items-center px-3 text-sm font-medium text-gray-500 hover:text-gray-900 border-b-2 border-transparent transition-colors"
										>
											{item.name}
										</a>
									)
								}
									return (
										<Link
											key={item.name}
											to={item.href}
											aria-current={isCurrent ? "page" : undefined}
											className={classNames(
												isCurrent
													? "border-b-2 border-[#55F] text-[#55F] font-semibold"
													: "border-b-2 border-transparent text-gray-500 hover:text-gray-900 hover:border-gray-300",
												"inline-flex items-center px-3 text-sm font-medium transition-colors"
											)}
										>
											{item.name}
										</Link>
									)
								})}
							</div>
						</div>
					</div>
					<div className="hidden sm:flex items-center sm:ml-6">
						{user ? (
							<Menu as="div" className="relative ml-3">
								<MenuButton className="relative flex items-center justify-center rounded-full ring-2 ring-transparent hover:ring-gray-200 transition-all">
									<span className="absolute -inset-1.5" />
									<span className="sr-only">Open user menu</span>
									<img alt="" src="user.png" className="size-8 rounded-full" />
								</MenuButton>

								<MenuItems transition className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-lg border border-gray-200 bg-white py-1 shadow-lg">
									<MenuItem>
										<Link to="/profile" className="block px-4 py-2 text-sm text-gray-700 hover:bg-gray-50">
											Your profile
										</Link>
									</MenuItem>
									<div className="my-1 border-t border-gray-100" />
									<MenuItem>
										<button
											type="button"
											onClick={handleLogout}
											className="block w-full px-4 py-2 text-left text-sm text-red-600 hover:bg-red-50"
										>
											Sign out
										</button>
									</MenuItem>
								</MenuItems>
							</Menu>
						) : (
							<div className="flex items-center gap-3">
								<Link to="/login" className="text-sm font-medium text-gray-600 hover:text-gray-900 transition-colors">
									Login
								</Link>
								<Link
									to="/register"
									className="rounded-full border border-[#55F] px-4 py-1.5 text-sm font-semibold text-[#55F] hover:bg-[#55F] hover:text-white transition-colors"
								>
									Register
								</Link>
							</div>
						)}
					</div>
				</div>
			</div>

			<DisclosurePanel className="sm:hidden border-t border-gray-100">
				<div className="space-y-0.5 px-4 pb-3 pt-2">
					{navigation.map((item) => {
						const isCurrent = !item.external && (item.href === "/"
							? location.pathname === item.href
							: location.pathname.startsWith(item.href))
						if (item.external) {
							return (
								<a
									key={item.name}
									href={item.href}
									target="_blank"
									rel="noopener noreferrer"
									className="block rounded-md px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900"
								>
									{item.name}
								</a>
							)
						}
						return (
							<DisclosureButton
								key={item.name}
								as={Link}
								to={item.href}
								aria-current={isCurrent ? "page" : undefined}
								className={classNames(
									isCurrent
										? "bg-indigo-50 text-[#55F] font-semibold"
										: "text-gray-600 hover:bg-gray-50 hover:text-gray-900",
									"block rounded-md px-3 py-2 text-sm font-medium"
								)}
							>
								{item.name}
							</DisclosureButton>
						)
					})}
					<div className="border-t border-gray-200 mt-2 pt-2">
						{user ? (
							<>
								<DisclosureButton
									as={Link}
									to="/profile"
									className="block rounded-md px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900"
								>
									Your profile
								</DisclosureButton>
								<DisclosureButton
									as="button"
									type="button"
									onClick={handleLogout}
									className="block w-full text-left rounded-md px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50"
								>
									Sign out
								</DisclosureButton>
							</>
						) : (
							<>
								<DisclosureButton
									as={Link}
									to="/login"
									className="block rounded-md px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900"
								>
									Login
								</DisclosureButton>
								<DisclosureButton
									as={Link}
									to="/register"
									className="block rounded-md px-3 py-2 text-sm font-semibold text-[#55F] hover:bg-indigo-50"
								>
									Register
								</DisclosureButton>
							</>
						)}
					</div>
				</div>
			</DisclosurePanel>
		</Disclosure>
	)
}
