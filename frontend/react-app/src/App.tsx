import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Todo } from './components/Todo'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Todo />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App