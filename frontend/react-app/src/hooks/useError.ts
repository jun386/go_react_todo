import { useNavigate } from 'react-router-dom'
import useStore from '../store'

export const useError = () => {
  const navigate = useNavigate()
  const resetEditedTask = useStore((state) => state.resetEditedTask)
  const switchErrorHandling = (msg: string) => {
    switch (msg) {
      default:
        alert(msg)
    }
  }
  return { switchErrorHandling }
}