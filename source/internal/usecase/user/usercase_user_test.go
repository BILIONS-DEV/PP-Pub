package user

import (
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
	"source/internal/repo/user/mocks"
	"testing"
)

type inputLoginSchema struct {
	inp  LoginInput
	want loginWant
}

type loginWant struct {
	err string
}

func TestUser_Login(t *testing.T) {

	userRepoMocks := new(mocks.RepoUser)

	repos := &repo.Repositories{User: userRepoMocks}
	userUsecase := NewUserFeUsecase(repos, nil)

	//=> cái này sử dụng trong trường hợp mình truyền fiber.Ctx vào usecase Method: `t.useCases.User.Login(email, password, ctx)`
	// Ví dụ như trường hợp mình truyền fiber.Ctx vào usecase để set cookie login.
	// tuy nhiên là không nên dùng như vậy vì hàm đó sẽ không sử dụng được ở worker hoặc 1 client nào khác mà k phải fiber
	//app := fiber.New()
	//ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	//_, err := userUsecase.Login(&input, ctx)

	/**
	tham khảo:
	https://blog.lamida.org/mocking-in-golang-using-testify/
	*/
	//=> nếu muốn kiểm tra từng trường hợp (chạy từ trên xuống dưới) vs đối số (Argument) đúng như yêu cầu (khai báo ở hàm `.On`) thì setup như dưới đây.
	//userRepoMocks.On("FindByLogin", "tungdt@83.com.vn", "password").Return(model.User{ID: 1, Email: "tungdt@83.com.vn", Permission: model.UserPermissionAdmin}, nil)
	//userRepoMocks.On("FindByLogin", "tungdtdev@gmail.com.vn", "password2").Return(model.User{}, errors.New(`record not found`))
	//userRepoMocks.On("FindByLogin", "tungdtdev2@gmail.com.vn", "password3").Return(model.User{}, errors.New(`what the error`))

	//=> còn nếu trường hợp k cần quan tâm đối số là gì & 1 mock.On cho tất cả các trường hợp thì viết như dưới đây:
	//userRepoMocks.On("FindByLogin", mock.Anything, mock.Anything).Return(model.User{ID: 1, Email: "tungdt@83.com.vn", Permission: model.UserPermissionAdmin}, nil)

	//=> còn trường hợp k quan tâm đối số, nhưng mỗi case là 1 mock khác nhau thì dùng thêm hàm .Once
	//=> mặc định là nếu dùng full đối số = `mock.Anything` thì hàm đó sẽ được coi là hàm dùng đi dùng lại, nếu muốn nó chạy 1 lần thì phải thêm .Once()
	// còn nếu mình nhập đối số cụ thể thì hàm đó là hàm dùng 1 lần.
	userRepoMocks.On("FindByLogin", mock.Anything, mock.Anything).Return(model.User{ID: 1, Email: "tungdt@83.com.vn", Permission: model.UserPermissionAdmin}, nil).Once()
	userRepoMocks.On("FindByLogin", mock.Anything, "password2").Return(model.User{}, errors.New(`record not found`))
	userRepoMocks.On("FindByLogin", mock.Anything, mock.Anything).Return(model.User{}, errors.New(`what the error`))

	inputs := []inputLoginSchema{
		{
			inp:  LoginInput{Email: "tungdt@83.com.vn", Password: "password"},
			want: loginWant{err: ""},
		},
		{
			inp:  LoginInput{Email: "tungdtdev@gmail.com.vn", Password: "password2"},
			want: loginWant{err: "record not found"},
		},
		{
			inp:  LoginInput{Email: "tungdtdev2@gmail.com.vn", Password: "password2"}, //=> case này vì có `password="password2"` nên vẫn thỏa mãn hàm On("FindByLogin", mock.Anything, "password2") nên nó nhận mock ở hàm đó.
			want: loginWant{err: "record not found"},
		},
		{
			inp:  LoginInput{Email: "tungdtdev3@gmail.com.vn", Password: "password3"},
			want: loginWant{err: "what the error"},
		},
	}

	assertions := require.New(t)
	for _, inp := range inputs {
		token, err := userUsecase.Login(&inp.inp)
		if err != nil {
			assertions.EqualError(err, inp.want.err, `message gì đó để đánh dấu`)
		} else {
			assertions.NotEqual(token, "")
		}
		fmt.Printf("\n token of email `%s`: %+v \n", inp.inp.Email, token)
	}

}
