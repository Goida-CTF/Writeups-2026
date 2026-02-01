package initial

import (
	"os"

	log "github.com/sirupsen/logrus"

	"vibecoded/internal/usecases"
)

var flag = os.Getenv("FLAG")

var initialUsers = []InitialUser{
	{
		Username: "elliot",
		Password: "Mr.Moejoe73!",
		IsAdmin:  true,
		Notes: []InitialNote{
			{
				"Шифрование",
				"Когда кент дал списать и ты не понял его подчерк",
			},
			{
				"Файрвол",
				"Когда друг пришел тебе бить лицо, но у тебя закрыта дверь в квартиру",
			},
			{
				"OSINT",
				"Когда узнал свой адрес по открытому журналу на столе классного руководителя",
			},
			{
				"Open Source",
				"Когда открыл тетрадь и показал ответы",
			},
			{
				"Утечка данных",
				"Когда ты хотел срать, но не успел добежать до толкана",
			},
			{
				"DDoS атака",
				"Срал под дверь соседу с разной одеждой",
			},
			{
				"Анонимность",
				"Когда не написал Ф.И.О. на листке с контрольной",
			},
			{
				"Пинг",
				"Когда решил проспать первую пару и проснулся к последней",
			},
			{
				"Облачное хранилище",
				"Когда в школе спратял шпору под партой",
			},
			{
				"Спуфинг",
				"Когда на контрольной подписал работу именем друга и ответил везде неправильно",
			},
			{
				"Ботнет",
				"Договорились всей группой не идти на первую пару",
			},
			{
				"Брутфорс",
				"На тесте написал ответы наугад и получил хорошую оценку " + flag,
			},
			{
				"VPN",
				"Когда передал шпору своему кенту через друга",
			},
			{
				"Кетфишинг",
				"Когда скачал войсмод, поставил женский голос и позвонил другу",
			},
			{
				"Бэкап",
				"Когда сохранил фотку друга после бара",
			},
			{
				"Реверс инжиниринг",
				"Когда пытаешься разобрать подчерк кента",
			},
			{
				"Сжатие файла",
				"Когда сократил каждое слово в шпоре до одной буквы",
			},
			{
				"Редирект",
				"Когда зашел в толкан, а там занято",
			},
		},
	},
	{
		Username: "admin2",
		Password: `jYR*Ne%N9wjoks_zVHCf!ZsBetq!b!Gf`,
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"README",
				"LOL u're 100% reading the writeup rn, there is no way...",
			},
		},
	},
	{
		Username: "john",
		Password: "!.YaNeThElIzOnDo.!4121",
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"Анекдот",
				"Тут должен был быть анекдот, but i forgor 💀",
			},
		},
	},
	{
		Username: "Marisa",
		Password: "BabyGirl#1",
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"Флаг",
				"🇷🇺",
			},
		},
	},
	{
		Username: "gopher",
		Password: "MNgopher42$@",
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"The goroutine joke",
				"Why did the goroutine cross the road?\n\nTo concurrently fetch the CTF flag, of course! 😂😂😂",
			},
		},
	},
	{
		Username: "alabaster1996",
		Password: "MyWIFE10ve$$ME",
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"Review",
				"Ну вы конечно круто напрограммировали",
			},
		},
	},
	{
		Username: "amogus",
		Password: "LoLLipop123!",
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"sus",
				"When the imposter is sus... and also hiding the CTF flag 😳",
			},
		},
	},
	{
		Username: "zombie",
		Password: "P@$$w0rd",
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"My fav quotes of all time 🗣🔥🔥🔥",
				`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nunc tincidunt turpis est, sed ornare tortor aliquet non. Nulla vel consectetur libero. Nulla volutpat, tellus vitae accumsan sagittis, felis mauris laoreet nibh, sit amet consequat lacus augue sed nunc. Aliquam in ex eu libero molestie pharetra. Duis a ligula augue. Vivamus sem dui, aliquet quis elementum ac, tristique sit amet nulla. Vivamus sagittis ultrices ligula eget lobortis.`,
			},
		},
	},
	{
		Username: "test",
		Password: "tEST.1234",
		IsAdmin:  false,
		Notes: []InitialNote{
			{
				"test",
				"test",
			},
		},
	},
}

func Initialize(uc *usecases.UseCases) {
	init := newInit(uc)

	for _, user := range initialUsers[:] {
		if err := init.createInitialUser(&user); err != nil {
			log.Fatalln("init.createInitialUser:", err)
		}
	}
}
