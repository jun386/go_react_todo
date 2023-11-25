import { FormEvent } from 'react'
import useStore from '../store'
import { useQueryTasks } from '../hooks/useQueryTasks'
import { useMutateTask } from '../hooks/useMutateTask'
import { TaskItem } from './TaskItem'

export const Todo = () => {
  const { editedTask } = useStore()
  const updateTask = useStore((state) => state.updateEditedTask)
  const { data, isLoading } = useQueryTasks()
  const { createTaskMutation, updateTaskMutation } = useMutateTask()
  const submitTaskHandler = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    if (editedTask.id === 0)
      createTaskMutation.mutate({
        title: editedTask.title,
      })
    else {
      updateTaskMutation.mutate(editedTask)
    }
  }
  return (
    <div className="flex justify-center items-center flex-col min-h-screen text-gray-600 font-mono">
      <div className="flex items-center my-3">
        <span className="text-center text-3xl font-extrabold">
          TODOリスト
        </span>
      </div>
      <form onSubmit={submitTaskHandler}>
        <input
          className="mb-3 mr-3 px-3 py-2 border border-gray-300"
          placeholder="タイトル"
          type="text"
          onChange={(e) => updateTask({ ...editedTask, title: e.target.value })}
          value={editedTask.title || ''}
        />
        <button
          className="disabled:opacity-40 mx-3 py-2 px-3 text-white bg-indigo-600 rounded"
          disabled={!editedTask.title}
        >
          {editedTask.id === 0 ? '作成' : '更新'}
        </button>
      </form>
      {isLoading ? (
        <p>ロード中...</p>
      ) : (
        <ul className="my-5">
          {data?.map((task) => (
            <TaskItem key={task.id} id={task.id} title={task.title} />
          ))}
        </ul>
      )}
    </div>
  )
}