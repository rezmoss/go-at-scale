// Example 45
// Testable service with dependencies
type UserService struct {
    repo   Repository
    hasher PasswordHasher
    mailer EmailSender
}

// Test doubles
type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Find(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*User), args.Error(1)
}

// Table-driven tests
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   User
        mockFn  func(*MockRepository)
        wantErr bool
    }{
        {
            name: "successful creation",
            input: User{Email: "test@example.com"},
            mockFn: func(repo *MockRepository) {
                repo.On("Save", mock.Anything, mock.AnythingOfType("*User")).
                    Return(nil)
            },
            wantErr: false,
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := new(MockRepository)
            tt.mockFn(repo)
            
            service := NewUserService(repo)
            err := service.CreateUser(context.Background(), tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
            }
            
            repo.AssertExpectations(t)
        })
    }
}