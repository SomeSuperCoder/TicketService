package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/SomeSuperCoder/OnlineShop/internal"
	"github.com/SomeSuperCoder/OnlineShop/internal/embeddings"
	"github.com/SomeSuperCoder/OnlineShop/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
)

func main() {
	ctx := context.Background()

	// Load configuration
	config := internal.LoadAppConfig()

	// Connect to database
	pool, repo, _ := internal.DatabaseConnect(ctx, config)
	defer pool.Close()

	log.Println("🚀 Starting mock data generation...")

	// Clean existing data
	if err := cleanData(ctx, pool); err != nil {
		log.Fatalf("Failed to clean data: %v", err)
	}

	// Generate mock data
	if err := generateMockData(ctx, repo, pool); err != nil {
		log.Fatalf("Failed to generate mock data: %v", err)
	}

	log.Println("✅ Mock data generation completed successfully!")
	printSummary(ctx, repo)
}

func cleanData(ctx context.Context, pool *pgxpool.Pool) error {
	log.Println("🧹 Cleaning existing data...")

	queries := []string{
		"TRUNCATE TABLE ticket_history CASCADE",
		"TRUNCATE TABLE ticket_tags CASCADE",
		"TRUNCATE TABLE comments CASCADE",
		"TRUNCATE TABLE complaint_details CASCADE",
		"TRUNCATE TABLE tickets CASCADE",
		"TRUNCATE TABLE users CASCADE",
		"TRUNCATE TABLE subcategories CASCADE",
		"TRUNCATE TABLE categories CASCADE",
		"TRUNCATE TABLE departments CASCADE",
		"TRUNCATE TABLE tags CASCADE",
		"TRUNCATE TABLE sources CASCADE",
		"TRUNCATE TABLE statuses CASCADE",
		"ALTER SEQUENCE categories_id_seq RESTART WITH 1",
		"ALTER SEQUENCE subcategories_id_seq RESTART WITH 1",
		"ALTER SEQUENCE departments_id_seq RESTART WITH 1",
		"ALTER SEQUENCE tags_id_seq RESTART WITH 1",
		"ALTER SEQUENCE sources_id_seq RESTART WITH 1",
		"ALTER SEQUENCE statuses_id_seq RESTART WITH 1",
	}

	for _, query := range queries {
		if _, err := pool.Exec(ctx, query); err != nil {
			return fmt.Errorf("failed to execute %s: %w", query, err)
		}
	}

	return nil
}

func generateMockData(ctx context.Context, repo *repository.Queries, pool *pgxpool.Pool) error {
	// 1. Create dictionaries
	log.Println("📚 Creating dictionaries...")
	if err := createDictionaries(ctx, repo); err != nil {
		return err
	}

	// 2. Create users
	log.Println("👥 Creating users...")
	if err := createUsers(ctx, repo); err != nil {
		return err
	}

	// 3. Create tickets with real embeddings
	log.Println("🎫 Creating tickets with AI embeddings...")
	if err := createTickets(ctx, repo); err != nil {
		return err
	}

	// 4. Add comments
	log.Println("💬 Adding comments...")
	if err := createComments(ctx, repo); err != nil {
		return err
	}

	// 5. Add tags
	log.Println("🏷️  Adding tags...")
	if err := createTags(ctx, repo); err != nil {
		return err
	}

	return nil
}

func createDictionaries(ctx context.Context, repo *repository.Queries) error {
	// Categories
	categories := []string{
		"ЖКХ",
		"Дороги и транспорт",
		"Благоустройство",
		"Безопасность и правопорядок",
		"Связь и телевидение",
		"Здравоохранение",
		"Образование",
		"Экология",
	}

	for _, name := range categories {
		if _, err := repo.CreateCategory(ctx, repository.CreateCategoryParams{Name: name}); err != nil {
			return err
		}
	}

	// Subcategories
	subcategories := map[int32][]string{
		1: {"Отопление", "Водоснабжение", "Электроснабжение", "Газоснабжение", "Плата за услуги"},
		2: {"Ремонт дорог", "Уборка дорог", "Освещение дорог", "Общественный транспорт", "Парковки"},
		3: {"Уборка территории", "Детские площадки", "Озеленение", "Ремонт подъездов", "Скверы и парки"},
		4: {"Бродячие животные", "Освещение улиц", "Видеонаблюдение", "Шум и беспорядок"},
		5: {"Качество интернета", "Мобильная связь", "Цифровое ТВ"},
		6: {"Поликлиники", "Больницы", "Скорая помощь", "Аптеки"},
		7: {"Детские сады", "Школы", "Дополнительное образование"},
		8: {"Мусор и свалки", "Загрязнение воздуха", "Загрязнение воды", "Шум"},
	}

	for catID, subs := range subcategories {
		for _, name := range subs {
			if _, err := repo.CreateSubcategory(ctx, repository.CreateSubcategoryParams{
				CategoryID: catID,
				Name:       name,
			}); err != nil {
				return err
			}
		}
	}

	// Departments
	departments := []string{
		"Администрация города Чебоксары",
		"Отдел ЖКХ",
		"Дорожный отдел",
		"Отдел благоустройства",
		"Отдел безопасности",
		"Департамент здравоохранения",
		"Департамент образования",
		"Экологический отдел",
	}

	for _, name := range departments {
		if _, err := repo.CreateDepartment(ctx, repository.CreateDepartmentParams{Name: name}); err != nil {
			return err
		}
	}

	// Tags
	tags := []string{
		"Срочно",
		"Требует проверки",
		"Дубликат",
		"Важное",
		"Массовая проблема",
		"Требует согласования",
		"Ожидает финансирования",
	}

	for _, name := range tags {
		if _, err := repo.CreateTag(ctx, repository.CreateTagParams{Name: name}); err != nil {
			return err
		}
	}

	// Sources
	sources := []string{
		"Веб-форма",
		"Госуслуги",
		"Горячая линия",
		"Мобильное приложение",
		"Социальные сети",
		"Электронная почта",
	}

	for _, name := range sources {
		if _, err := repo.CreateSource(ctx, repository.CreateSourceParams{Name: name}); err != nil {
			return err
		}
	}

	return nil
}

func createUsers(ctx context.Context, repo *repository.Queries) error {
	users := []struct {
		id           string
		email        string
		role         repository.UserRole
		status       repository.UserStatus
		departmentID *int32
		firstName    string
		lastName     string
		middleName   string
	}{
		// Admin
		{"11111111-1111-1111-1111-111111111111", "admin@cheboksary.ru", repository.UserRoleAdmin, repository.UserStatusActive, ptr(int32(1)), "Иван", "Иванов", "Иванович"},

		// Organizations
		{"22222222-2222-2222-2222-222222222222", "roi.gkh@cheboksary.ru", repository.UserRoleOrg, repository.UserStatusActive, ptr(int32(2)), "Петр", "Петров", "Петрович"},
		{"33333333-3333-3333-3333-333333333333", "roi.roads@cheboksary.ru", repository.UserRoleOrg, repository.UserStatusActive, ptr(int32(3)), "Сергей", "Сергеев", "Сергеевич"},
		{"44444444-4444-4444-4444-444444444444", "roi.improvement@cheboksary.ru", repository.UserRoleOrg, repository.UserStatusActive, ptr(int32(4)), "Мария", "Смирнова", "Александровна"},
		{"99999999-9999-9999-9999-999999999999", "roi.health@cheboksary.ru", repository.UserRoleOrg, repository.UserStatusBlocked, ptr(int32(6)), "Анна", "Антонова", "Андреевна"},

		// Executors
		{"55555555-5555-5555-5555-555555555555", "executor1@cheboksary.ru", repository.UserRoleExecutor, repository.UserStatusActive, ptr(int32(2)), "Алексей", "Алексеев", "Алексеевич"},
		{"66666666-6666-6666-6666-666666666666", "executor2@cheboksary.ru", repository.UserRoleExecutor, repository.UserStatusActive, ptr(int32(2)), "Ольга", "Орлова", "Олеговна"},
		{"77777777-7777-7777-7777-777777777777", "executor3@cheboksary.ru", repository.UserRoleExecutor, repository.UserStatusActive, ptr(int32(3)), "Дмитрий", "Дмитриев", "Дмитриевич"},
		{"88888888-8888-8888-8888-888888888888", "executor4@cheboksary.ru", repository.UserRoleExecutor, repository.UserStatusActive, ptr(int32(4)), "Елена", "Егорова", "Евгеньевна"},
		{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", "executor5@cheboksary.ru", repository.UserRoleExecutor, repository.UserStatusBlocked, ptr(int32(5)), "Николай", "Николаев", "Николаевич"},
	}

	for _, u := range users {
		id, _ := uuid.Parse(u.id)
		if _, err := repo.CreateUser(ctx, repository.CreateUserParams{
			ID:           id,
			Email:        u.email,
			Role:         u.role,
			Status:       u.status,
			DepartmentID: u.departmentID,
			FirstName:    &u.firstName,
			LastName:     &u.lastName,
			MiddleName:   &u.middleName,
		}); err != nil {
			return err
		}
	}

	return nil
}

func ptr[T any](v T) *T {
	return &v
}

type TicketData struct {
	description   string
	subcategoryID int32
	departmentID  *int32
	status        repository.TicketStatus
	createdAt     time.Time
	senderName    string
	senderPhone   *string
	senderEmail   *string
	lat           float64
	lng           float64
}

func createTickets(ctx context.Context, repo *repository.Queries) error {
	// Diverse and interesting ticket scenarios
	tickets := []TicketData{
		// CLOSED TICKETS (resolved, 15-45 days old)
		{
			description:   "В нашем доме на улице Ленина 45 уже неделю нет отопления. Температура в квартирах опустилась до 15 градусов. Жильцы, особенно дети и пожилые люди, мерзнут. Просим срочно решить проблему!",
			subcategoryID: 1, status: repository.TicketStatusClosed, departmentID: ptr(int32(2)),
			createdAt:  time.Now().Add(-45 * 24 * time.Hour),
			senderName: "Иванов Иван Иванович", senderPhone: ptr("+7 900 123 45 67"), senderEmail: ptr("ivanov@example.com"),
			lat: 56.1324, lng: 47.2501,
		},
		{
			description:   "На перекрестке улиц Гагарина и Мира образовалась огромная яма глубиной около 30 см. Несколько машин уже повредили колеса. Требуется срочный ремонт дорожного покрытия.",
			subcategoryID: 6, status: repository.TicketStatusClosed, departmentID: ptr(int32(3)),
			createdAt:  time.Now().Add(-40 * 24 * time.Hour),
			senderName: "Петрова Мария Сергеевна", senderPhone: ptr("+7 900 234 56 78"), senderEmail: nil,
			lat: 56.1340, lng: 47.2520,
		},
		{
			description:   "Уличное освещение на улице Победы не работает уже месяц. Вечером и ночью очень темно и опасно ходить. Были случаи краж и нападений. Просим восстановить освещение.",
			subcategoryID: 12, status: repository.TicketStatusClosed, departmentID: ptr(int32(5)),
			createdAt:  time.Now().Add(-35 * 24 * time.Hour),
			senderName: "Сидоров Петр Алексеевич", senderPhone: nil, senderEmail: ptr("sidorov@example.com"),
			lat: 56.1310, lng: 47.2480,
		},
		{
			description:   "Во дворе дома №12 по улице Калинина не убирается мусор. Контейнеры переполнены, мусор валяется вокруг. Появились крысы. Ситуация антисанитарная!",
			subcategoryID: 11, status: repository.TicketStatusClosed, departmentID: ptr(int32(4)),
			createdAt:  time.Now().Add(-30 * 24 * time.Hour),
			senderName: "Анонимный житель", senderPhone: ptr("+7 900 345 67 89"), senderEmail: nil,
			lat: 56.1330, lng: 47.2510,
		},
		{
			description:   "Детская площадка возле дома №8 на улице Чапаева в аварийном состоянии. Качели сломаны, горка проржавела, песочница грязная. Дети могут травмироваться!",
			subcategoryID: 12, status: repository.TicketStatusClosed, departmentID: ptr(int32(4)),
			createdAt:  time.Now().Add(-25 * 24 * time.Hour),
			senderName: "Козлова Ольга Викторовна", senderPhone: ptr("+7 900 456 78 90"), senderEmail: ptr("kozlova@example.com"),
			lat: 56.1320, lng: 47.2490,
		},
		{
			description:   "В поликлинике №3 огромные очереди к терапевту. Запись на прием только через 2 недели. Люди вынуждены стоять с утра, чтобы попасть к врачу. Нужны дополнительные специалисты!",
			subcategoryID: 21, status: repository.TicketStatusClosed, departmentID: ptr(int32(6)),
			createdAt:  time.Now().Add(-20 * 24 * time.Hour),
			senderName: "Морозова Анна Петровна", senderPhone: ptr("+7 900 567 89 01"), senderEmail: ptr("morozova@example.com"),
			lat: 56.1350, lng: 47.2530,
		},

		// OPEN TICKETS (in progress, 2-6 days old)
		{
			description:   "Холодные батареи в доме на улице Ленина 67. Отопление работает, но температура в квартирах всего 18 градусов. Возможно, проблема с циркуляцией или давлением в системе.",
			subcategoryID: 1, status: repository.TicketStatusOpen, departmentID: ptr(int32(2)),
			createdAt:  time.Now().Add(-6 * 24 * time.Hour),
			senderName: "Смирнов Алексей Иванович", senderPhone: ptr("+7 900 678 90 12"), senderEmail: ptr("smirnov@example.com"),
			lat: 56.1325, lng: 47.2505,
		},
		{
			description:   "Из крана течет ржавая вода с неприятным запахом. Невозможно использовать для питья и приготовления пищи. Проблема появилась 3 дня назад после ремонтных работ на водопроводе.",
			subcategoryID: 2, status: repository.TicketStatusOpen, departmentID: ptr(int32(2)),
			createdAt:  time.Now().Add(-5 * 24 * time.Hour),
			senderName: "Новикова Елена Сергеевна", senderPhone: nil, senderEmail: ptr("novikova@example.com"),
			lat: 56.1335, lng: 47.2515,
		},
		{
			description:   "Дорога на улице Мира полностью разбита. Асфальт в ямах и трещинах по всей длине улицы. Проезд затруднен, машины объезжают по встречной полосе. Очень опасно!",
			subcategoryID: 6, status: repository.TicketStatusOpen, departmentID: ptr(int32(3)),
			createdAt:  time.Now().Add(-4 * 24 * time.Hour),
			senderName: "Волков Дмитрий Александрович", senderPhone: ptr("+7 900 789 01 23"), senderEmail: nil,
			lat: 56.1345, lng: 47.2525,
		},
		{
			description:   "После снегопада прошло 3 дня, а во дворе дома №15 по улице Советской снег так и не убрали. Невозможно пройти к подъезду, машины застревают. Пожилым людям особенно тяжело.",
			subcategoryID: 7, status: repository.TicketStatusOpen, departmentID: ptr(int32(3)),
			createdAt:  time.Now().Add(-3 * 24 * time.Hour),
			senderName: "Лебедев Сергей Николаевич", senderPhone: ptr("+7 900 890 12 34"), senderEmail: ptr("lebedev@example.com"),
			lat: 56.1315, lng: 47.2485,
		},
		{
			description:   "В нашем районе нужна новая детская площадка. Ближайшая находится в 15 минутах ходьбы. У нас много семей с детьми, а играть негде. Есть подходящее место возле дома №20.",
			subcategoryID: 12, status: repository.TicketStatusOpen, departmentID: ptr(int32(4)),
			createdAt:  time.Now().Add(-2 * 24 * time.Hour),
			senderName: "Соколова Татьяна Владимировна", senderPhone: ptr("+7 900 901 23 45"), senderEmail: nil,
			lat: 56.1305, lng: 47.2495,
		},
		{
			description:   "В школе №12 протекает крыша в спортзале. После дождя на полу лужи, занятия физкультурой приходится отменять. Дети лишены возможности заниматься спортом.",
			subcategoryID: 26, status: repository.TicketStatusOpen, departmentID: ptr(int32(7)),
			createdAt:  time.Now().Add(-24 * time.Hour),
			senderName: "Григорьева Ирина Петровна", senderPhone: ptr("+7 900 012 34 56"), senderEmail: ptr("grigorieva@example.com"),
			lat: 56.1360, lng: 47.2540,
		},

		// OVERDUE TICKETS (open for 8-20 days - will show in overdue stats)
		{
			description:   "Протечка крыши в подъезде дома №5 на улице Пушкина. Вода течет по стенам, образовалась плесень. Штукатурка отваливается. Опасно для здоровья жильцов!",
			subcategoryID: 14, status: repository.TicketStatusOpen, departmentID: ptr(int32(2)),
			createdAt:  time.Now().Add(-20 * 24 * time.Hour),
			senderName: "Федорова Наталья Ивановна", senderPhone: ptr("+7 900 123 45 67"), senderEmail: ptr("fedorova@example.com"),
			lat: 56.1300, lng: 47.2500,
		},
		{
			description:   "Огромная яма на дороге по улице Строителей глубиной около 40 см и шириной 2 метра. Уже было несколько аварий. Машины повреждают подвеску. Ситуация критическая!",
			subcategoryID: 6, status: repository.TicketStatusOpen, departmentID: ptr(int32(3)),
			createdAt:  time.Now().Add(-18 * 24 * time.Hour),
			senderName: "Михайлов Андрей Сергеевич", senderPhone: nil, senderEmail: ptr("mikhailov@example.com"),
			lat: 56.1350, lng: 47.2530,
		},
		{
			description:   "Стая бродячих собак (около 10 особей) обосновалась во дворе дома №7 по улице Комсомольской. Собаки агрессивные, нападают на людей. Дети боятся выходить гулять. Требуется отлов!",
			subcategoryID: 15, status: repository.TicketStatusOpen, departmentID: ptr(int32(5)),
			createdAt:  time.Now().Add(-15 * 24 * time.Hour),
			senderName: "Павлова Ирина Александровна", senderPhone: ptr("+7 900 234 56 78"), senderEmail: nil,
			lat: 56.1295, lng: 47.2475,
		},
		{
			description:   "В парке Победы совершенно нет освещения. Вечером парк превращается в опасную зону. Были случаи грабежей и нападений. Люди боятся ходить через парк после захода солнца.",
			subcategoryID: 16, status: repository.TicketStatusOpen, departmentID: ptr(int32(5)),
			createdAt:  time.Now().Add(-12 * 24 * time.Hour),
			senderName: "Романов Виктор Петрович", senderPhone: ptr("+7 900 345 67 89"), senderEmail: ptr("romanov@example.com"),
			lat: 56.1360, lng: 47.2540,
		},
		{
			description:   "Интернет в нашем районе (улицы Садовая, Цветочная, Зеленая) работает крайне плохо. Скорость падает до 1 Мбит/с, постоянные обрывы связи. Невозможно работать удаленно!",
			subcategoryID: 18, status: repository.TicketStatusOpen, departmentID: nil,
			createdAt:  time.Now().Add(-10 * 24 * time.Hour),
			senderName: "Кузнецова Людмила Николаевна", senderPhone: nil, senderEmail: ptr("kuznetsova@example.com"),
			lat: 56.1290, lng: 47.2470,
		},
		{
			description:   "Возле дома №25 на улице Лесной образовалась несанкционированная свалка. Люди выбрасывают строительный мусор, старую мебель, бытовую технику. Свалка растет, появился неприятный запах.",
			subcategoryID: 29, status: repository.TicketStatusOpen, departmentID: ptr(int32(8)),
			createdAt:  time.Now().Add(-9 * 24 * time.Hour),
			senderName: "Николаев Олег Владимирович", senderPhone: ptr("+7 900 456 78 90"), senderEmail: nil,
			lat: 56.1285, lng: 47.2465,
		},

		// INIT TICKETS (new, not assigned, recent)
		{
			description:   "Счета за ЖКУ выросли в 2 раза за последние 3 месяца без объяснения причин. Тарифы вроде не менялись, потребление тоже. Просим разобраться в ситуации и пересчитать.",
			subcategoryID: 5, status: repository.TicketStatusInit, departmentID: nil,
			createdAt:  time.Now().Add(-5 * time.Hour),
			senderName: "Егоров Максим Андреевич", senderPhone: ptr("+7 900 567 89 01"), senderEmail: ptr("egorov@example.com"),
			lat: 56.1328, lng: 47.2512,
		},
		{
			description:   "Предлагаю организовать озеленение двора дома №18 на улице Молодежной. Сейчас там только асфальт и парковка. Хотелось бы посадить деревья, кустарники, разбить клумбы.",
			subcategoryID: 13, status: repository.TicketStatusInit, departmentID: nil,
			createdAt:  time.Now().Add(-3 * time.Hour),
			senderName: "Белова Светлана Игоревна", senderPhone: ptr("+7 900 678 90 12"), senderEmail: nil,
			lat: 56.1318, lng: 47.2488,
		},
		{
			description:   "Мобильная связь оператора МТС в районе улицы Промышленной практически не работает. Звонки обрываются, интернет не грузится. Проблема появилась после установки нового оборудования.",
			subcategoryID: 19, status: repository.TicketStatusInit, departmentID: nil,
			createdAt:  time.Now().Add(-2 * time.Hour),
			senderName: "Зайцев Артем Дмитриевич", senderPhone: nil, senderEmail: ptr("zaitsev@example.com"),
			lat: 56.1342, lng: 47.2522,
		},
		{
			description:   "Автобус №15 ходит с большими задержками. По расписанию должен быть каждые 15 минут, но на самом деле интервал 30-40 минут. Люди опаздывают на работу.",
			subcategoryID: 9, status: repository.TicketStatusInit, departmentID: nil,
			createdAt:  time.Now().Add(-1 * time.Hour),
			senderName: "Орлов Игорь Сергеевич", senderPhone: ptr("+7 900 789 01 23"), senderEmail: ptr("orlov@example.com"),
			lat: 56.1308, lng: 47.2498,
		},
		{
			description:   "В детском саду №8 не хватает мест. Очередь на зачисление составляет более 200 человек. Дети 3-4 лет не могут попасть в садик. Родители вынуждены сидеть дома или искать частные варианты.",
			subcategoryID: 25, status: repository.TicketStatusInit, departmentID: nil,
			createdAt:  time.Now().Add(-30 * time.Minute),
			senderName: "Васильева Екатерина Александровна", senderPhone: ptr("+7 900 890 12 34"), senderEmail: nil,
			lat: 56.1355, lng: 47.2535,
		},
	}

	for i, t := range tickets {
		log.Printf("  Creating ticket %d/%d: %s...", i+1, len(tickets), t.description[:50])

		// Generate real embedding
		embedding, err := embeddings.GetEmbedding(t.description)
		if err != nil {
			log.Printf("    ⚠️  Warning: Failed to generate embedding, using fallback: %v", err)
			// Create a simple fallback embedding
			embedding = createFallbackEmbedding()
		}

		// Create ticket
		ticket, err := repo.CreateTicketWithDefaults(ctx, repository.CreateTicketWithDefaultsParams{
			Description:   t.description,
			SubcategoryID: t.subcategoryID,
			DepartmentID:  t.departmentID,
			Embedding:     embedding,
		})
		if err != nil {
			return fmt.Errorf("failed to create ticket: %w", err)
		}

		// Update status and created_at
		if t.status != repository.TicketStatusInit {
			if _, err := repo.UpdateTicketSimple(ctx, repository.UpdateTicketSimpleParams{
				ID: ticket.ID,
				Status: repository.NullTicketStatus{
					TicketStatus: t.status,
					Valid:        true,
				},
			}); err != nil {
				return err
			}
		}

		// Update created_at (direct SQL since sqlc doesn't have this query)
		if err := repo.UpdateTicketCreatedAt(ctx, repository.UpdateTicketCreatedAtParams{
			ID:        ticket.ID,
			CreatedAt: t.createdAt,
		}); err != nil {
			return err
		}

		// Create complaint details
		geoLocation := fmt.Sprintf("POINT(%f %f)", t.lng, t.lat)
		if _, err := repo.CreateComplaint(ctx, repository.CreateComplaintParams{
			Ticket:      ticket.ID,
			Description: t.description,
			SenderName:  t.senderName,
			SenderPhone: t.senderPhone,
			SenderEmail: t.senderEmail,
			GeoLocation: geoLocation,
		}); err != nil {
			return err
		}

		// Create history entry
		newValue := fmt.Sprintf(`{"status": "%s", "subcategory_id": %d}`, t.status, t.subcategoryID)
		if _, err := repo.CreateHistoryEntry(ctx, repository.CreateHistoryEntryParams{
			TicketID: ticket.ID,
			Action:   repository.HistoryActionCreated,
			NewValue: []byte(newValue),
		}); err != nil {
			return err
		}

		// Add status change history for non-init tickets
		if t.status != repository.TicketStatusInit {
			oldValue := `{"status": "init"}`
			newValue := fmt.Sprintf(`{"status": "%s"}`, t.status)
			if _, err := repo.CreateHistoryEntryWithTime(ctx, repository.CreateHistoryEntryWithTimeParams{
				TicketID:  ticket.ID,
				Action:    repository.HistoryActionStatusChanged,
				OldValue:  []byte(oldValue),
				NewValue:  []byte(newValue),
				CreatedAt: t.createdAt.Add(1 * time.Hour),
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func createFallbackEmbedding() *pgvector.Vector {
	// Create a simple random embedding as fallback
	emb := make([]float32, 768)
	for i := range emb {
		emb[i] = float32(i) / 768.0
	}
	v := pgvector.NewVector(emb)
	return &v
}

func createComments(ctx context.Context, repo *repository.Queries) error {
	// Get all tickets
	tickets, err := repo.ListAllTickets(ctx)
	if err != nil {
		return err
	}

	comments := map[string][]string{
		// Closed tickets have resolution comments
		"closed": {
			"Заявка принята в работу",
			"Специалист выехал на место",
			"Проблема устранена, работы завершены",
		},
		// Open tickets have work in progress
		"open": {
			"Заявка зарегистрирована",
			"Передано в профильный отдел",
			"Работы запланированы",
		},
	}

	for _, ticket := range tickets {
		var commentList []string
		switch ticket.Status {
		case repository.TicketStatusClosed:
			commentList = comments["closed"]
		case repository.TicketStatusOpen:
			commentList = comments["open"]
		default:
			continue
		}

		for _, msg := range commentList {
			if _, err := repo.CreateComment(ctx, repository.CreateCommentParams{
				Ticket:  ticket.ID,
				Message: msg,
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func createTags(ctx context.Context, repo *repository.Queries) error {
	// Get all tickets
	tickets, err := repo.ListAllTickets(ctx)
	if err != nil {
		return err
	}

	// Tag assignment logic
	for _, ticket := range tickets {
		var tags []int32

		// Overdue tickets get "Срочно" tag
		if ticket.Status == repository.TicketStatusOpen {
			age := time.Since(ticket.CreatedAt)
			if age > 7*24*time.Hour {
				tags = append(tags, 1) // Срочно
			}
			if age > 10*24*time.Hour {
				tags = append(tags, 4) // Важное
			}
		}

		// Some tickets need verification
		if ticket.SubcategoryID == 5 || ticket.SubcategoryID == 18 {
			tags = append(tags, 2) // Требует проверки
		}

		// Mass problems
		if ticket.SubcategoryID == 15 || ticket.SubcategoryID == 29 {
			tags = append(tags, 5) // Массовая проблема
		}

		// Add tags
		for _, tagID := range tags {
			if _, err := repo.AddTagsToTicket(ctx, repository.AddTagsToTicketParams{
				Ticket: ticket.ID,
				Tags:   []int32{tagID},
			}); err != nil {
				// Ignore duplicate errors
				continue
			}
		}
	}

	return nil
}

func printSummary(ctx context.Context, repo *repository.Queries) {
	log.Println("\n" + strings.Repeat("=", 50))
	log.Println("📊 MOCK DATA SUMMARY")
	log.Println(strings.Repeat("=", 50))

	// Count categories
	categories, _ := repo.GetCategories(ctx)
	log.Printf("Categories: %d", len(categories))

	// Count departments
	departments, _ := repo.GetDepartments(ctx)
	log.Printf("Departments: %d", len(departments))

	// Count users
	users, _ := repo.ListUsers(ctx, repository.ListUsersParams{Limit: 1000})
	log.Printf("Users: %d", len(users))

	// Count tickets by status
	allTickets, _ := repo.ListAllTickets(ctx)
	var closed, open, init int
	for _, t := range allTickets {
		switch t.Status {
		case repository.TicketStatusClosed:
			closed++
		case repository.TicketStatusOpen:
			open++
		case repository.TicketStatusInit:
			init++
		}
	}

	log.Printf("Tickets: %d", len(allTickets))
	log.Printf("  - Closed: %d", closed)
	log.Printf("  - Open: %d", open)
	log.Printf("  - Init: %d", init)

	// Count overdue
	overdue, _ := repo.GetOverdueTickets(ctx, repository.GetOverdueTicketsParams{
		MinLostDays: 0,
		Limit:       1000,
	})
	log.Printf("  - Overdue: %d", len(overdue))

	// Count comments
	commentCount := 0
	for _, t := range allTickets {
		comments, _ := repo.GetCommentsForTicket(ctx, repository.GetCommentsForTicketParams{Ticket: t.ID})
		commentCount += len(comments)
	}
	log.Printf("Comments: %d", commentCount)

	log.Println(strings.Repeat("=", 50))
	log.Println("✨ Ready for frontend integration!")
	log.Println(strings.Repeat("=", 50))
}
