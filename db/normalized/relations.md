## Описание переменных отношений
* ### profile
  Персональные данные пользователя.

* ### tag
  Множество всех возможных тегов, уникальных по title. Предполагается только пополнение данного множества.

* ### pin
  Пины (картинки). Пины могут быть приватными, которые видны только автору, или публичные - доступные всем для просмотра.
* ### pin_tag
  Отношение служащее связью M:M для [pin](#pin) и [tag](#tag).
* ### like_pin
  Оценки "`нравится`" пользователей к пинами.
* ### board
  Доски (коллекция [пинов](#pin)). Доски как и пины могут быть приватными, которые видны только автору, или публичные - доступные всем для просмотра.
* ### board_tag
  Отношение служащее связью M:M для [board](#board) и [tag](#tag).
* ### subscription_board
  Информация о подписках пользователя на доску.
* ### membership
  Информация о размещении пина на доске.
* ### role
  Роли для соавторов досок. Предполагается две: `read-write` - просмотр и добавление содержимого доски, `read-only` - только просмотр содержимого, приобретает смысл для приватных досок.
* ### contributor
  Соавторы досок со своими ролями.
* ### subscription_user
  Информация о подписках пользователя на другого пользователя.
* ### comment
  Комментарии к пинам.
* ### like_comment
  Оценки "`нравится`" пользователей к комментариям.
* ### message
  Сообщения личных переписок пользователей.

Часть сущностей имеют дополнительный атрибут `deleted_at`, несущий следующую функциональность: при логическом удалении сущности данный атрибут обновляется до времени когда это произошло, сама запись из БД не удаляется.

## Описание функциональных зависимостей
#### Relation [profile](#profile):
{id} -> {username, email, password, avatar, name, surname, about_me, created_at, updated_at, deleted_at}
{username} -> {id, email, password, avatar, name, surname, about_me, created_at, updated_at, deleted_at}
{email} -> {id, username, password, avatar, name, surname, about_me, created_at, updated_at, deleted_at}

#### Relation [tag](#tag):
{id} -> {title, created_at}
{title} -> {id, created_at}

#### Relation [pin](#pin):
{id} -> {author, title, description, picture, public, created_at, updated_at, deleted_at}

#### Relation [pin_tag](#pin_tag):
{pin_id, tag_id} -> {created_at}

#### Relation [like_pin](#like_pin):
{pin_id, user_id} -> {created_at}

#### Relation [board](#board):
{id} -> {author, title, description, public, created_at, updated_at, deleted_at}

#### Relation [board_tag](#board_tag):
{board_id, tag_id} -> {created_at}

#### Relation [subscription_board](#subscription_board):
{board_id, user_id} -> {created_at}

#### Relation [membership](#membership):
{board_id, pin_id} -> {added_at}

#### Relation [role](#role):
{id} -> {name}
{name} -> {id}

#### Relation [contributor](#contributor):
{user_id, board_id} -> {role_id, added_at, updated_at}

#### Relation [subscription_user](#subscription_user):
{whom, who} -> {created_at}

#### Relation [comment](#comment):
{id} -> {author, pin_id, content, created_at, updated_at, deleted_at}

#### Relation [like_comment](#like_comment):
{comment_id, user_id} -> {created_at}

#### Relation [message](#message):
{id} -> {user_from, user_to, content, created_at, updated_at, deleted_at}


## Соответствие нормальной форме
```
1-ая НФ - Переменная отношения находится в первой нормальной форме (1НФ) тогда и только тогда, когда в любом допустимом значении отношения каждый его кортеж содержит только одно значение для каждого из атрибутов.
```

Как видно у всех сущностей атрибуты представляют собой единственное значение простого типа, следовательно 1 НФ выполняется.
#
```
2-ая НФ - Переменная отношения находится во второй нормальной форме тогда и только тогда, когда она находится в первой нормальной форме, и каждый неключевой атрибут неприводимо (функционально полно) зависит от ее потенциального ключа.
```

Так как 1 НФ уже выполнена и функциональные зависимости от части ключа отсутсвуют, то 2 НФ тоже выполняется.
#
```
3-я НФ - Переменная отношения находится в третьей нормальной форме, когда она находится во второй нормальной форме, и отсутствуют транзитивные функциональные зависимости неключевых атрибутов от ключевых.
```
Так как 2 НФ уже выполнена и каждый атрибут переменной отношения зависит только от ключевого атрибута, то 3 НФ тоже выполнятся.
#
```
НФ Бойса-Кодда - Отношение находится в НФБК, когда каждая нетривиальная и неприводимая слева функциональная зависимость обладает потенциальным ключом в качестве детерминанта.
```
Как видно все детерминанты являются потенциальными ключами, следовательно НФБК выполняется.
