### Микросервис в этом репозитории

В данном репозитории находится код, который предназначен для работы *offline-магазина* в **pet-проекте**.

#### Описание проекта

Проект будет имитировать бэкэнд сети магазинов по продаже часов ⌚. Описание требований:
- продажи будут вестись как в обычных *offline-магазинах*, так и через интернет магазин
- возможна отправка товара не только со складов компании, но и из *offline-магазинов*
- доступна бронь товара через интернет для самовывоза из конкретного *offline-магазина*  

#### Ограничения, принятые в проекте:
+ номер заказа - положительное число
+ номер заказа через интернет должен быть больше десяти. Номера меньше десяти используются для касс, которые отмечают
пробитые, но еще не оплаченные товары, как заказанные. Это необходимо для того, чтобы исключить бронирование товаров, 
находящихся в процессе продажи
+ Дефекты товаров или упаковки, влияющие на цену, шифруются в артикуле товара (а это значит, что товар без дефектов и с 
дефектом имеют разные артикулы). Дефекты кодируются следующим образом - за основу берется артикул неповрежденного
товара, после ставится точка, а далее идут четыре цифры:
  1. 0 - корпус без повреждений, 1 - корпус имеет легкие царапины, 2 - корпус имеет сильные царапины
  2. 0 - дисплей/стекло без повреждений, 1 - дисплей/стекло имеет легкие царапины, 2 - дисплей/стекло имеет сильные
  царапины. (Для ремешков всегда указывается 0, так как у них нет дисплея)
  3. 0 - упаковка/коробка не вскрывалась, 1 - упаковка/коробка вскрывалась
  4. 0 - упаковка/коробка без повреждений, 1 - упаковка/коробка повреждена

#### Для чего это всё написано?

В данном репозитории содержится код, являющийся частью моего **pet-проекта**, цель которого - изучение языка Golang,
микросервисной архитектуры, взаимодействия микросервисов между собой. Возможно, содержимое этого и имеющих к нему
отношения репозиториев поможет кому-нибудь определиться, достаточный ли у меня уровень для взаимного сотрудничества 😉