package todos

type Repositorier interface {
	GetTodos(*[]Todo) error
	NewTodo(*Todo) error
	DeleteTodo(*Todo, int) error
}
