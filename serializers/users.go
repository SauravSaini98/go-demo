package serializers

import (
    "my_project/models"
)

type UserSerializer struct {
    user models.User
}

func NewUserSerializer(user models.User) UserSerializer {
    return UserSerializer{user: user}
}

func (us UserSerializer) Serialize() map[string]interface{} {
    return map[string]interface{}{
        "id":    us.user.ID,
        "email": us.user.Email,
        "name": us.user.Name,
    }
}
