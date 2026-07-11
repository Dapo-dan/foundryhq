// A React "component" is just a function that returns JSX (the HTML-like
// syntax below). React calls this function to figure out what to display,
// and re-calls it whenever the component's state/props change. Component
// function names must start with a capital letter — that's how React (and
// JSX) tells your own components apart from plain HTML tags like <div>.
function App() {
  return (
    // JSX looks like HTML but it's actually JavaScript. Key differences:
    // - `className` instead of `class` (since `class` is a reserved JS word)
    // - the strings here (e.g. "flex min-h-svh items-center justify-center")
    //   are Tailwind CSS utility classes — each one maps to a small CSS rule
    //   (flex = display:flex, min-h-svh = min-height:100svh, etc.) instead of
    //   writing a separate .css file.
    <div className="flex min-h-svh items-center justify-center">
      <h1 className="text-2xl font-medium">FoundryHQ</h1>
    </div>
  )
}

// main.tsx imports this as `App` and renders it — this is the root/top-level
// component that everything else in the app will eventually nest inside.
export default App
