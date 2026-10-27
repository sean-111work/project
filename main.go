package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// User 结构体定义
type User struct {
	ID   int
	Name string
	Age  int
}

// UserStore 封装用户数据及状态，避免全局变量带来的并发与测试隔离问题
type UserStore struct {
	users  []User
	nextID int
}

// NewUserStore 初始化 UserStore
func NewUserStore() *UserStore {
	return &UserStore{
		users:  make([]User, 0),
		nextID: 1,
	}
}

// AddUser 添加用户方法
func (s *UserStore) AddUser(name string, age int) int {
	id := s.nextID
	newUser := User{
		ID:   id,
		Name: name,
		Age:  age,
	}
	s.users = append(s.users, newUser)
	s.nextID++
	return id
}

// FindUserByID 根据 ID 查找用户 (当前为 O(n) 线性扫描，未来可升级为 map[int]User 实现 O(1) 查找)
func (s *UserStore) FindUserByID(id int) (User, bool) {
	for _, user := range s.users {
		if user.ID == id {
			return user, true
		}
	}
	return User{}, false
}

// GetAllUsers 获取所有用户列表
func (s *UserStore) GetAllUsers() []User {
	return s.users
}

// readInput 安全读取输入，返回 (string, error) 以支持 EOF 和异常流处理
func readInput(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		// 检查是否有读取错误，若无错误则是 EOF（如管道输入结束、Ctrl+D等）
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return "", io.EOF
	}
	return strings.TrimSpace(scanner.Text()), nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	store := NewUserStore()

	// 主循环：菜单持续运行
	for {
		fmt.Println("\n============================")
		fmt.Println("    用户管理系统 (User Mgmt)")
		fmt.Println("============================")
		fmt.Println("1. 添加用户")
		fmt.Println("2. 根据 ID 查询用户")
		fmt.Println("3. 展示所有用户")
		fmt.Println("0. 退出系统")
		fmt.Print("请选择操作 (1-3, 0退出): ")

		choice, err := readInput(scanner)
		if err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("\n检测到输入流结束 (EOF)，安全退出系统。")
				break
			}
			fmt.Printf("\n【错误】读取输入发生异常: %v\n", err)
			break
		}

		// switch 语句根据输入执行对应功能
		switch choice {
		case "1":
			fmt.Println("\n--- 添加用户 ---")
			fmt.Print("请输入姓名: ")
			name, err := readInput(scanner)
			if err != nil {
				fmt.Println("【错误】读取姓名失败。")
				continue
			}
			if name == "" {
				fmt.Println("【错误】姓名不能为空！")
				continue
			}

			fmt.Print("请输入年龄: ")
			ageStr, err := readInput(scanner)
			if err != nil {
				fmt.Println("【错误】读取年龄失败。")
				continue
			}

			age, err := strconv.Atoi(ageStr)
			if err != nil || age <= 0 || age > 150 {
				fmt.Println("【错误】无效的年龄输入！请输入 1 到 150 之间的整数。")
				continue
			}

			id := store.AddUser(name, age)
			fmt.Printf("【成功】用户添加成功！分配的 ID 为: %d\n", id)

		case "2":
			fmt.Println("\n--- 查询用户 ---")
			users := store.GetAllUsers()
			if len(users) == 0 {
				fmt.Println("【提示】当前用户列表为空。")
				continue
			}

			fmt.Print("请输入要查询的用户 ID: ")
			idStr, err := readInput(scanner)
			if err != nil {
				fmt.Println("【错误】读取 ID 失败。")
				continue
			}

			id, err := strconv.Atoi(idStr)
			if err != nil {
				fmt.Println("【错误】请输入有效的数字 ID！")
				continue
			}

			if user, found := store.FindUserByID(id); found {
				fmt.Printf("【找到用户】ID: %d | 姓名: %s | 年龄: %d\n", user.ID, user.Name, user.Age)
			} else {
				fmt.Printf("【错误】未找到 ID 为 %d 的用户！\n", id)
			}

		case "3":
			fmt.Println("\n--- 用户列表 ---")
			users := store.GetAllUsers()
			if len(users) == 0 {
				fmt.Println("【提示】暂无任何用户数据。")
				return
			}

			fmt.Printf("%-5s %-15s %-5s\n", "ID", "姓名", "年龄")
			fmt.Println("-----------------------------")
			for _, user := range users {
				fmt.Printf("%-5d %-15s %-5d\n", user.ID, user.Name, user.Age)
			}

		case "0", "4":
			fmt.Println("感谢使用，系统已退出！")
			return

		default:
			fmt.Println("【错误】无效的选择，请输入数字 0 - 3！")
		}
	}
}
