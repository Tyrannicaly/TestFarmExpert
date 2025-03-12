Из документации Peppol UBL Invoice так: <cbc:LineExtensionAmount> — чистая сумма строк, <cac:TaxTotal> — общая сумма налогов, <cac:TaxSubtotal> — разбивка налогов.

Моя логика
Парсинг: суммы (например, "1500.00") в int64 через big.Float (150000 центов), проценты (например, "20") в int через strconv.Atoi.
Валидация: для каждого <cac:TaxSubtotal> проверяется TaxableAmount * Percent / 100 = TaxAmount (например, 100000 * 20 / 100 = 20000).
Сумма налогов: сумма TaxAmount из <cac:TaxSubtotal> (25000) должна равняться <cbc:TaxAmount> в <cac:TaxTotal>.
Итог: LineExtensionAmount + TaxAmount (например, 150000 + 35000 = 185000).


Для точности использую int64 избегая число с плавающей точки (float64).Добавил юнит тесты чтобы убедиться что код корректный.
