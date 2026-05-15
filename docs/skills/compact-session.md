# Skill: Compact Session — รับมือเมื่อ Context Window เต็ม

เมื่อ context window ของ AI session เต็ม ให้ทำ **Compact Session** เพื่อสรุปสถานะงานปัจจุบัน
แล้วเปิด session ใหม่ต่อจากจุดที่ค้างไว้ โดยไม่สูญเสียความต่อเนื่อง

---

## สัญญาณที่บ่งบอกว่า Context Window ใกล้เต็ม

- AI ตอบช้าลง หรือ ตอบซ้ำสิ่งที่ทำไปแล้ว
- AI "ลืม" สิ่งที่คุยไปตั้งแต่ต้น session
- ได้รับข้อความแจ้งเตือนว่า context ยาวเกินไป
- AI ไม่ reference โค้ดหรือไฟล์ที่เพิ่งแก้ไปก่อนหน้า

---

## ขั้นตอน Compact Session

### ขั้นที่ 1: สรุปสถานะปัจจุบัน (Summarize)

ขอให้ AI สรุปงานที่ทำอยู่ โดยพิมพ์ prompt นี้:

```
สรุปสถานะงานปัจจุบันให้ฉันหน่อย ได้แก่:
1. งานที่ทำเสร็จไปแล้ว
2. งานที่กำลังทำอยู่ (in-progress)
3. งานที่ยังเหลืออยู่
4. ไฟล์ที่แก้ไขไปแล้ว
5. context สำคัญที่ต้องรู้ก่อนทำงานต่อ
```

### ขั้นที่ 2: บันทึกสรุปไว้ใน session memory

ขอให้ AI บันทึกสรุปนั้นลง memory ก่อนปิด session:

```
บันทึกสรุปนี้ลง session memory ก่อนเลย
```

AI จะใช้ `/memories/session/` เพื่อเก็บ plan และ context ที่จำเป็น

### ขั้นที่ 3: เปิด Session ใหม่

เปิด chat session ใหม่ แล้วเริ่มต้นด้วย prompt นี้:

```
อ่าน session memory ก่อน แล้วทำงานต่อจากที่ค้างไว้
```

---

## Template: Handoff Prompt สำหรับ Session ใหม่

```
Context สำคัญ:
- โปรเจกต์: [ชื่อโปรเจกต์]
- ภาษา/Framework: Go + Gin + GORM + PostgreSQL
- งานที่กำลังทำ: [อธิบายสั้นๆ]
- ไฟล์ที่แก้ล่าสุด: [รายชื่อไฟล์]
- สิ่งที่ทำเสร็จไปแล้ว:
  - [x] item 1
  - [x] item 2
- สิ่งที่ยังเหลือ:
  - [ ] item 3
  - [ ] item 4
- ปัญหาที่พบ (ถ้ามี): [อธิบาย]

ขั้นตอนต่อไป: [บอกว่าต้องทำอะไรก่อน]
```

---

## ตัวอย่างการใช้งานจริง

สมมติกำลังแก้ไข unit test ใน expensesController_test.go ค้างอยู่:

```
Context สำคัญ:
- โปรเจกต์: lab-golang (Go + Gin + GORM + PostgreSQL)
- งานที่กำลังทำ: เขียน unit test สำหรับ expenses controller
- ไฟล์ที่แก้ล่าสุด: test/controllers/expensesController_test.go

- สิ่งที่ทำเสร็จ:
  - [x] TestCreateExpenses_Success
  - [x] TestCreateExpenses_InvalidJSON

- สิ่งที่ยังเหลือ:
  - [ ] TestGetListMoneyCard_Success
  - [ ] TestGetListMoneyCard_Empty
  - [ ] TestDeleteExpensesDetail_Success

ขั้นตอนต่อไป: เขียน TestGetListMoneyCard_Success โดยใช้ Table-Driven Test pattern
```

---

## Best Practices

| แนวทาง | รายละเอียด |
|---|---|
| บันทึก memory บ่อยๆ | อย่ารอให้ context เต็มค่อยบันทึก ควรทำเป็นระยะ |
| ระบุไฟล์ที่แก้ | ระบุชื่อไฟล์และบรรทัดสำคัญเสมอ เพื่อให้ session ใหม่ pick up ต่อได้ทันที |
| TODO list | เก็บ todo list ไว้ใน memory หรือในโค้ด เช่น `// TODO:` |
| Commit บ่อยๆ | commit โค้ดก่อน compact เพื่อไม่ให้งานหาย |
| ระบุ pattern ที่ตัดสินใจไว้ | เช่น "ตัดสินใจใช้ Table-Driven Test" เพื่อไม่ให้ session ใหม่เปลี่ยนแนวทาง |

---

## Quick Cheatsheet

```
# เมื่อ context ใกล้เต็ม
1. พิมพ์: "สรุปสถานะงานและบันทึกลง session memory"
2. commit โค้ดที่ทำเสร็จแล้ว: git add . && git commit -m "wip: ..."
3. เปิด session ใหม่
4. พิมพ์: "อ่าน session memory แล้วทำงานต่อจากที่ค้างไว้"
```
