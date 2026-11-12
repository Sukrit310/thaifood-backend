-- Idempotent reference data and 30 real restaurants.
-- Sources and individual photo licenses: docs/RESTAURANT_DATA_SOURCES.md
-- Retrieved 2026-11-12. Ratings are calculated only from this site's reviews.
BEGIN;

INSERT INTO provinces(name_th,region) VALUES
('กรุงเทพมหานคร','ภาคกลาง'),
('กำแพงเพชร','ภาคกลาง'),
('ชัยนาท','ภาคกลาง'),
('นครนายก','ภาคกลาง'),
('นครปฐม','ภาคกลาง'),
('นครสวรรค์','ภาคกลาง'),
('นนทบุรี','ภาคกลาง'),
('ปทุมธานี','ภาคกลาง'),
('พระนครศรีอยุธยา','ภาคกลาง'),
('พิจิตร','ภาคกลาง'),
('พิษณุโลก','ภาคกลาง'),
('เพชรบูรณ์','ภาคกลาง'),
('ลพบุรี','ภาคกลาง'),
('สมุทรปราการ','ภาคกลาง'),
('สมุทรสงคราม','ภาคกลาง'),
('สมุทรสาคร','ภาคกลาง'),
('สระบุรี','ภาคกลาง'),
('สิงห์บุรี','ภาคกลาง'),
('สุโขทัย','ภาคกลาง'),
('สุพรรณบุรี','ภาคกลาง'),
('อ่างทอง','ภาคกลาง'),
('อุทัยธานี','ภาคกลาง'),
('เชียงราย','ภาคเหนือ'),
('เชียงใหม่','ภาคเหนือ'),
('น่าน','ภาคเหนือ'),
('พะเยา','ภาคเหนือ'),
('แพร่','ภาคเหนือ'),
('แม่ฮ่องสอน','ภาคเหนือ'),
('ลำปาง','ภาคเหนือ'),
('ลำพูน','ภาคเหนือ'),
('อุตรดิตถ์','ภาคเหนือ'),
('กาฬสินธุ์','ภาคตะวันออกเฉียงเหนือ'),
('ขอนแก่น','ภาคตะวันออกเฉียงเหนือ'),
('ชัยภูมิ','ภาคตะวันออกเฉียงเหนือ'),
('นครพนม','ภาคตะวันออกเฉียงเหนือ'),
('นครราชสีมา','ภาคตะวันออกเฉียงเหนือ'),
('บึงกาฬ','ภาคตะวันออกเฉียงเหนือ'),
('บุรีรัมย์','ภาคตะวันออกเฉียงเหนือ'),
('มหาสารคาม','ภาคตะวันออกเฉียงเหนือ'),
('มุกดาหาร','ภาคตะวันออกเฉียงเหนือ'),
('ยโสธร','ภาคตะวันออกเฉียงเหนือ'),
('ร้อยเอ็ด','ภาคตะวันออกเฉียงเหนือ'),
('เลย','ภาคตะวันออกเฉียงเหนือ'),
('ศรีสะเกษ','ภาคตะวันออกเฉียงเหนือ'),
('สกลนคร','ภาคตะวันออกเฉียงเหนือ'),
('สุรินทร์','ภาคตะวันออกเฉียงเหนือ'),
('หนองคาย','ภาคตะวันออกเฉียงเหนือ'),
('หนองบัวลำภู','ภาคตะวันออกเฉียงเหนือ'),
('อำนาจเจริญ','ภาคตะวันออกเฉียงเหนือ'),
('อุดรธานี','ภาคตะวันออกเฉียงเหนือ'),
('อุบลราชธานี','ภาคตะวันออกเฉียงเหนือ'),
('จันทบุรี','ภาคตะวันออก'),
('ฉะเชิงเทรา','ภาคตะวันออก'),
('ชลบุรี','ภาคตะวันออก'),
('ตราด','ภาคตะวันออก'),
('ปราจีนบุรี','ภาคตะวันออก'),
('ระยอง','ภาคตะวันออก'),
('สระแก้ว','ภาคตะวันออก'),
('กาญจนบุรี','ภาคตะวันตก'),
('ตาก','ภาคตะวันตก'),
('ประจวบคีรีขันธ์','ภาคตะวันตก'),
('เพชรบุรี','ภาคตะวันตก'),
('ราชบุรี','ภาคตะวันตก'),
('กระบี่','ภาคใต้'),
('ชุมพร','ภาคใต้'),
('ตรัง','ภาคใต้'),
('นครศรีธรรมราช','ภาคใต้'),
('นราธิวาส','ภาคใต้'),
('ปัตตานี','ภาคใต้'),
('พังงา','ภาคใต้'),
('พัทลุง','ภาคใต้'),
('ภูเก็ต','ภาคใต้'),
('ยะลา','ภาคใต้'),
('ระนอง','ภาคใต้'),
('สงขลา','ภาคใต้'),
('สตูล','ภาคใต้'),
('สุราษฎร์ธานี','ภาคใต้')
ON CONFLICT(name_th) DO NOTHING;

INSERT INTO categories(name_th) VALUES
('อาหารเหนือ'),('อาหารอีสาน'),('อาหารใต้'),('อาหารไทย'),('คาเฟ่'),('อาหารทะเล'),('อาหารมังสวิรัติ')
ON CONFLICT(name_th) DO NOTHING;

-- Retire only the old, explicitly identified fictional seed rows. Protect
-- child writes while checking references, so a concurrent review/favorite
-- cannot be cascaded away between NOT EXISTS and DELETE. Real/admin rows
-- and fictional rows edited by an admin or referenced by a member are preserved.
SET LOCAL lock_timeout = '5s';
LOCK TABLE reviews, favorites IN SHARE MODE;
DELETE FROM restaurants r
WHERE r.source_name = 'ThaiFood Guide author-created fictional sample'
  AND r.source_url LIKE 'urn:thaifood-guide:fictional-sample:%'
  AND r.updated_at = r.created_at
  AND NOT EXISTS (SELECT 1 FROM reviews v WHERE v.restaurant_id = r.id)
  AND NOT EXISTS (SELECT 1 FROM favorites f WHERE f.restaurant_id = r.id);

WITH real_places(name,province_name,category_name,address,description,latitude,longitude,image_url,source_name,source_url) AS (VALUES
    ('เฮือนม่วนใจ๋',
     'เชียงใหม่',
     'อาหารเหนือ',
     '24 ราชพฤกษ์ ตำบลช้างเผือก อำเภอเมืองเชียงใหม่ จังหวัดเชียงใหม่ 50300',
     'ร้านอาหารเหนือในย่านช้างเผือก เมืองเชียงใหม่ เสิร์ฟอาหารพื้นเมืองล้านนา

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20241118154150001
ภาพ: Lanna cuisine starters.JPG — Takeaway / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=18431732 (แสดงเป็นภาพย่อ)',
     18.7998212,
     98.9754871,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/1/18/Lanna_cuisine_starters.JPG/960px-Lanna_cuisine_starters.JPG',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20241118154150001'),
    ('เมธาวลัย ศรแดง',
     'กรุงเทพมหานคร',
     'อาหารไทย',
     '78 2 แขวงบวรนิเวศ เขตพระนคร กรุงเทพมหานคร 10200',
     'ร้านอาหารไทยในเขตพระนคร กรุงเทพมหานคร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20200629140800464046
ภาพ: Methavalai Sorndaeng เมธาวลัย ศรแดง (Dec 2021).jpg — Chainwit. / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=114430590 (แสดงเป็นภาพย่อ)',
     13.756113,
     100.5020899,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/7/7a/Methavalai_Sorndaeng_%E0%B9%80%E0%B8%A1%E0%B8%98%E0%B8%B2%E0%B8%A7%E0%B8%A5%E0%B8%B1%E0%B8%A2_%E0%B8%A8%E0%B8%A3%E0%B9%81%E0%B8%94%E0%B8%87_%28Dec_2021%29.jpg/960px-Methavalai_Sorndaeng_%E0%B9%80%E0%B8%A1%E0%B8%98%E0%B8%B2%E0%B8%A7%E0%B8%A5%E0%B8%B1%E0%B8%A2_%E0%B8%A8%E0%B8%A3%E0%B9%81%E0%B8%94%E0%B8%87_%28Dec_2021%29.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20200629140800464046'),
    ('บลู เอเลเฟ่นท์ ภูเก็ต',
     'ภูเก็ต',
     'อาหารใต้',
     '96 ตำบลตลาดเหนือ อำเภอเมืองภูเก็ต จังหวัดภูเก็ต 83000',
     'ร้านอาหารไทยและอาหารพื้นเมืองภูเก็ตในอาคารสถาปัตยกรรมชิโนโปรตุกีส

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20250123132939001
ภาพ: Blue Elephant, Phuket, Thailand - 01.jpg — Sharon Hahn Darlin / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=115058205 (แสดงเป็นภาพย่อ)',
     7.886030099999998,
     98.38474409999999,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/1/14/Blue_Elephant%2C_Phuket%2C_Thailand_-_01.jpg/960px-Blue_Elephant%2C_Phuket%2C_Thailand_-_01.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20250123132939001'),
    ('ข้าวซอยแม่สาย',
     'เชียงใหม่',
     'อาหารเหนือ',
     'ตำบลช้างเผือก อำเภอเมืองเชียงใหม่ จังหวัดเชียงใหม่ 50300',
     'ร้านข้าวซอยและอาหารเหนือในย่านช้างเผือก เมืองเชียงใหม่

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20231116144610001
ภาพ: Khao soi kai.JPG — Takeaway / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=18291021 (แสดงเป็นภาพย่อ)',
     18.7995989,
     98.9752135,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/7/78/Khao_soi_kai.JPG/960px-Khao_soi_kai.JPG',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20231116144610001'),
    ('Breezeo – Exotic Beachfront Seafood',
     'ชลบุรี',
     'อาหารทะเล',
     'Royal Cliff Beach Hotel, 353 ถนนพระตำหนัก เมืองพัทยา จังหวัดชลบุรี 20150',
     'ห้องอาหารริมชายหาดของ Royal Cliff ในพัทยา ให้บริการอาหารทะเลและเมนูย่าง

ที่มาข้อมูล: Royal Cliff Hotels Group (https://www.royalcliff.com/restaurant/breezeo/)
ภาพ: 5.Breezeo - Be Yourself Dining.jpg — Cormaria / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=68347344 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/8/85/5.Breezeo_-_Be_Yourself_Dining.jpg/960px-5.Breezeo_-_Be_Yourself_Dining.jpg',
     'Royal Cliff Hotels Group',
     'https://www.royalcliff.com/restaurant/breezeo/'),
    ('เขาทองฮิลล์',
     'กระบี่',
     'คาเฟ่',
     '135 ตำบลเขาทอง อำเภอเมืองกระบี่ จังหวัดกระบี่ 81000',
     'คาเฟ่และร้านอาหารบนเขาในตำบลเขาทอง จังหวัดกระบี่ พร้อมทิวทัศน์ทะเล

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20220725162556741893
ภาพ: View from Khaothong Hill SF0001.jpg — Stefan Fussan / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=178144679 (แสดงเป็นภาพย่อ)',
     8.172732671876428,
     98.75266639547202,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/7/75/View_from_Khaothong_Hill_SF0001.jpg/960px-View_from_Khaothong_Hill_SF0001.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20220725162556741893'),
    ('ประจักษ์เป็ดย่าง',
     'กรุงเทพมหานคร',
     'อาหารไทย',
     '1415 เขตบางรัก กรุงเทพมหานคร 10500',
     'ร้านเป็ดย่างในเขตบางรัก กรุงเทพมหานคร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201202221118203228
ภาพ: เป็ดย่าง ประจักษ์.jpg — Chainwit. / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=114298367 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/1/11/%E0%B9%80%E0%B8%9B%E0%B9%87%E0%B8%94%E0%B8%A2%E0%B9%88%E0%B8%B2%E0%B8%87_%E0%B8%9B%E0%B8%A3%E0%B8%B0%E0%B8%88%E0%B8%B1%E0%B8%81%E0%B8%A9%E0%B9%8C.jpg/960px-%E0%B9%80%E0%B8%9B%E0%B9%87%E0%B8%94%E0%B8%A2%E0%B9%88%E0%B8%B2%E0%B8%87_%E0%B8%9B%E0%B8%A3%E0%B8%B0%E0%B8%88%E0%B8%B1%E0%B8%81%E0%B8%A9%E0%B9%8C.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201202221118203228'),
    ('Little Tree',
     'นครปฐม',
     'คาเฟ่',
     '43 หมู่ 5 อำเภอสามพราน จังหวัดนครปฐม 73110',
     'คาเฟ่และร้านอาหารในสวนริมแม่น้ำท่าจีน อำเภอสามพราน จังหวัดนครปฐม

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201305091129375151
ภาพ: The Little Tree Garden - panoramio.jpg — Piyatad / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=54179275 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/3/38/The_Little_Tree_Garden_-_panoramio.jpg/960px-The_Little_Tree_Garden_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201305091129375151'),
    ('หลู้ลำ',
     'เชียงราย',
     'อาหารเหนือ',
     '188/8 ม.20 อำเภอเมืองเชียงราย จังหวัดเชียงราย 57000',
     'ร้านอาหารพื้นเมืองเหนือในอำเภอเมืองเชียงราย

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201112161450361519
ภาพ: Sai Ua ไส้อั่ว at Lu Lam, Chiang Rai (2026) - IMG 01.jpg — Chainwit. / CC BY 4.0 (https://creativecommons.org/licenses/by/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=202181792 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/8/88/Sai_Ua_%E0%B9%84%E0%B8%AA%E0%B9%89%E0%B8%AD%E0%B8%B1%E0%B9%88%E0%B8%A7_at_Lu_Lam%2C_Chiang_Rai_%282026%29_-_IMG_01.jpg/960px-Sai_Ua_%E0%B9%84%E0%B8%AA%E0%B9%89%E0%B8%AD%E0%B8%B1%E0%B9%88%E0%B8%A7_at_Lu_Lam%2C_Chiang_Rai_%282026%29_-_IMG_01.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201112161450361519'),
    ('ข้าวเม่า-ข้าวฟ่าง',
     'เชียงใหม่',
     'อาหารเหนือ',
     '181 ม.7 ตำบลหนองควาย อำเภอหางดง จังหวัดเชียงใหม่',
     'ร้านอาหารไทยและอาหารเหนือในอำเภอหางดง จังหวัดเชียงใหม่

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201210161335504266
ภาพ: Khaoma o khaofang Restaurant (1) (12659492384).jpg — Iain Cameron / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=144946196 (แสดงเป็นภาพย่อ)',
     18.729942,
     98.94358260000001,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/b/b8/Khaoma_o_khaofang_Restaurant_%281%29_%2812659492384%29.jpg/960px-Khaoma_o_khaofang_Restaurant_%281%29_%2812659492384%29.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201210161335504266'),
    ('นายโซว',
     'กรุงเทพมหานคร',
     'อาหารไทย',
     '3/1 เขตป้อมปราบศัตรูพ่าย กรุงเทพมหานคร 10100',
     'ร้านอาหารในเขตป้อมปราบศัตรูพ่าย กรุงเทพมหานคร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201202231518163343
ภาพ: Nai Sow ร้านอาหาร นายโซว (2025) - img 01.jpg — Chainwit. / CC BY 4.0 (https://creativecommons.org/licenses/by/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=173726737 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/a/aa/Nai_Sow_%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%AD%E0%B8%B2%E0%B8%AB%E0%B8%B2%E0%B8%A3_%E0%B8%99%E0%B8%B2%E0%B8%A2%E0%B9%82%E0%B8%8B%E0%B8%A7_%282025%29_-_img_01.jpg/960px-Nai_Sow_%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%AD%E0%B8%B2%E0%B8%AB%E0%B8%B2%E0%B8%A3_%E0%B8%99%E0%B8%B2%E0%B8%A2%E0%B9%82%E0%B8%8B%E0%B8%A7_%282025%29_-_img_01.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201202231518163343'),
    ('ข้าวเปิ๊บยายเครื่อง',
     'สุโขทัย',
     'อาหารไทย',
     'อำเภอศรีสัชนาลัย จังหวัดสุโขทัย 64130',
     'ร้านข้าวเปิ๊บและอาหารท้องถิ่นในอำเภอศรีสัชนาลัย จังหวัดสุโขทัย

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20180216110359339866
ภาพ: Khao Poep "original" at Yai Krueang ข้าวเปิ๊บ ยายเครื่อง (2026) - img 01.jpg — Chainwit. / CC BY 4.0 (https://creativecommons.org/licenses/by/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=192576371 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/f/f2/Khao_Poep_%22original%22_at_Yai_Krueang_%E0%B8%82%E0%B9%89%E0%B8%B2%E0%B8%A7%E0%B9%80%E0%B8%9B%E0%B8%B4%E0%B9%8A%E0%B8%9A_%E0%B8%A2%E0%B8%B2%E0%B8%A2%E0%B9%80%E0%B8%84%E0%B8%A3%E0%B8%B7%E0%B9%88%E0%B8%AD%E0%B8%87_%282026%29_-_img_01.jpg/960px-Khao_Poep_%22original%22_at_Yai_Krueang_%E0%B8%82%E0%B9%89%E0%B8%B2%E0%B8%A7%E0%B9%80%E0%B8%9B%E0%B8%B4%E0%B9%8A%E0%B8%9A_%E0%B8%A2%E0%B8%B2%E0%B8%A2%E0%B9%80%E0%B8%84%E0%B8%A3%E0%B8%B7%E0%B9%88%E0%B8%AD%E0%B8%87_%282026%29_-_img_01.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20180216110359339866'),
    ('Blue Elephant',
     'กรุงเทพมหานคร',
     'อาหารไทย',
     '233 เขตสาทร กรุงเทพมหานคร 10120',
     'ร้านอาหารไทยและโรงเรียนสอนทำอาหารในอาคารเก่า ย่านสาทร กรุงเทพมหานคร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20200630112731464217
ภาพ: Blue Elephany cooking school Bangkom Sathon.jpg — Chainwit. / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=74702145 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/e/e5/Blue_Elephany_cooking_school_Bangkom_Sathon.jpg/960px-Blue_Elephany_cooking_school_Bangkom_Sathon.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20200630112731464217'),
    ('เเลเล กริลล์',
     'กระบี่',
     'อาหารทะเล',
     '83 ตำบลอ่าวนาง อำเภอเมืองกระบี่ จังหวัดกระบี่ 81180',
     'ร้านอาหารในย่านอ่าวนาง จังหวัดกระบี่ พร้อมทิวทัศน์ทะเล

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201201051108292258
ภาพ: Lae Lay Grill restaurant - view.jpg — kallerna / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=19122704 (แสดงเป็นภาพย่อ)',
     8.03831720945541,
     98.82165205296137,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/9/9e/Lae_Lay_Grill_restaurant_-_view.jpg/960px-Lae_Lay_Grill_restaurant_-_view.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201201051108292258'),
    ('บ้านใน',
     'กรุงเทพมหานคร',
     'อาหารไทย',
     '102/13 เขตพญาไท กรุงเทพมหานคร 10400',
     'ร้านอาหารไทยในบ้านใน เดอะเรมินิสเซนซ์ เขตพญาไท กรุงเทพมหานคร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20200630093301464090
ภาพ: Baannai Hotel, BKK, 2018-12-11 (01).jpg — 21janvier1793 / CC0 (https://creativecommons.org/publicdomain/zero/1.0/deed.en)
ที่มาภาพ: https://commons.wikimedia.org/?curid=110078785 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/d/d0/Baannai_Hotel%2C_BKK%2C_2018-12-11_%2801%29.jpg/960px-Baannai_Hotel%2C_BKK%2C_2018-12-11_%2801%29.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20200630093301464090'),
    ('โป่งแยงแอ่งดอย',
     'เชียงใหม่',
     'อาหารเหนือ',
     'เลขที่ 49/3 ม.2 ตำบลโป่งแยง อำเภอแม่ริม จังหวัดเชียงใหม่',
     'ร้านอาหารไทยและอาหารเหนือในตำบลโป่งแยง อำเภอแม่ริม จังหวัดเชียงใหม่

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201210161408204278
ภาพ: PongYaeng Ang Doi Restaurant (12661479483).jpg — Iain Cameron / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=144946448 (แสดงเป็นภาพย่อ)',
     18.8982297,
     98.8542932,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/0/06/PongYaeng_Ang_Doi_Restaurant_%2812661479483%29.jpg/960px-PongYaeng_Ang_Doi_Restaurant_%2812661479483%29.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201210161408204278'),
    ('ภัตตาคารจันทร์เพ็ญ',
     'กรุงเทพมหานคร',
     'อาหารไทย',
     '1020/1 เขตสาทร กรุงเทพมหานคร 10120',
     'ภัตตาคารอาหารไทยสาขาพระราม 4 กรุงเทพมหานคร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201202281311003650
ภาพ: Chandrphen Restaurant Rama 4 ภัตตาคารจันทร์เพ็ญ สาขาพระราม 4 (2026) - IMG 01.jpg — Chainwit. / CC BY 4.0 (https://creativecommons.org/licenses/by/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=199292210 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/2/2f/Chandrphen_Restaurant_Rama_4_%E0%B8%A0%E0%B8%B1%E0%B8%95%E0%B8%95%E0%B8%B2%E0%B8%84%E0%B8%B2%E0%B8%A3%E0%B8%88%E0%B8%B1%E0%B8%99%E0%B8%97%E0%B8%A3%E0%B9%8C%E0%B9%80%E0%B8%9E%E0%B9%87%E0%B8%8D_%E0%B8%AA%E0%B8%B2%E0%B8%82%E0%B8%B2%E0%B8%9E%E0%B8%A3%E0%B8%B0%E0%B8%A3%E0%B8%B2%E0%B8%A1_4_%282026%29_-_IMG_01.jpg/960px-Chandrphen_Restaurant_Rama_4_%E0%B8%A0%E0%B8%B1%E0%B8%95%E0%B8%95%E0%B8%B2%E0%B8%84%E0%B8%B2%E0%B8%A3%E0%B8%88%E0%B8%B1%E0%B8%99%E0%B8%97%E0%B8%A3%E0%B9%8C%E0%B9%80%E0%B8%9E%E0%B9%87%E0%B8%8D_%E0%B8%AA%E0%B8%B2%E0%B8%82%E0%B8%B2%E0%B8%9E%E0%B8%A3%E0%B8%B0%E0%B8%A3%E0%B8%B2%E0%B8%A1_4_%282026%29_-_IMG_01.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201202281311003650'),
    ('ร้านอาหารเกศริน',
     'ประจวบคีรีขันธ์',
     'อาหารทะเล',
     '17/1 อำเภอหัวหิน จังหวัดประจวบคีรีขันธ์ 77110',
     'ร้านอาหารทะเลในอำเภอหัวหิน จังหวัดประจวบคีรีขันธ์

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 20201012092230477068
ภาพ: Ketsarin Restaurant - panoramio.jpg — josef knecht / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=54382539 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/7/7a/Ketsarin_Restaurant_-_panoramio.jpg/960px-Ketsarin_Restaurant_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20201012092230477068'),
    ('คุ้มขันโตก',
     'เชียงใหม่',
     'อาหารเหนือ',
     '139 ม.4 ตำบลหนองป่าครั่ง อำเภอเมืองเชียงใหม่ จังหวัดเชียงใหม่ 50000',
     'ร้านอาหารเหนือพร้อมการแสดงที่คุ้มขันโตก จังหวัดเชียงใหม่

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201210161441344292
ภาพ: Khum Khantoke.jpg — whyyan / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=44800193 (แสดงเป็นภาพย่อ)',
     18.799491,
     99.0263702,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/d/d5/Khum_Khantoke.jpg/960px-Khum_Khantoke.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201210161441344292'),
    ('ประภาคาร ซีฟู้ด',
     'ชลบุรี',
     'อาหารทะเล',
     'บ้านอำเภอ อำเภอสัตหีบ จังหวัดชลบุรี',
     'ร้านอาหารทะเลย่านบ้านอำเภอ อำเภอสัตหีบ จังหวัดชลบุรี

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (https://data.go.th/th/dataset/restaurant1), PLACE_ID 201211131052574711
ภาพ: Prapakarn Restaurant - panoramio.jpg — Boonchai C / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=54350798 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/2/29/Prapakarn_Restaurant_-_panoramio.jpg/960px-Prapakarn_Restaurant_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201211131052574711')
)
INSERT INTO restaurants(name,province_id,category_id,address,description,latitude,longitude,image_url,source_name,source_url)
SELECT s.name,p.id,c.id,s.address,s.description,s.latitude::double precision,s.longitude::double precision,s.image_url,s.source_name,s.source_url
FROM real_places s
JOIN provinces p ON p.name_th=s.province_name
JOIN categories c ON c.name_th=s.category_name
-- Stable source IDs prevent duplicates and preserve later administrator edits.
ON CONFLICT(source_name,source_url) WHERE source_name IS NOT NULL AND source_url IS NOT NULL DO NOTHING;

-- Additional real places: ten previously empty provinces (2026-11-12).
-- INSERT-only: existing restaurants, reviews and favorites are preserved.
WITH additional_places(name,province_name,category_name,address,description,latitude,longitude,image_url,source_name,source_url) AS (VALUES
    ('รำพึง',
     'สิงห์บุรี',
     'อาหารไทย',
     'ตำบลบางมัญ อำเภอเมืองสิงห์บุรี จังหวัดสิงห์บุรี 16000',
     'ร้านอาหารไทยในอำเภอเมืองสิงห์บุรี จังหวัดสิงห์บุรี

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201112191540281739)
ภาพ: "รำพึง" ร้านอาหารอร่อยเมืองสิงห์ - panoramio.jpg — ::::=UT=:::: / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=52789509 (แสดงเป็นภาพย่อ)',
     14.8986018,
     100.4125322,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/d/d0/%22%E0%B8%A3%E0%B8%B3%E0%B8%9E%E0%B8%B6%E0%B8%87%22_%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%AD%E0%B8%B2%E0%B8%AB%E0%B8%B2%E0%B8%A3%E0%B8%AD%E0%B8%A3%E0%B9%88%E0%B8%AD%E0%B8%A2%E0%B9%80%E0%B8%A1%E0%B8%B7%E0%B8%AD%E0%B8%87%E0%B8%AA%E0%B8%B4%E0%B8%87%E0%B8%AB%E0%B9%8C_-_panoramio.jpg/960px-%22%E0%B8%A3%E0%B8%B3%E0%B8%9E%E0%B8%B6%E0%B8%87%22_%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%AD%E0%B8%B2%E0%B8%AB%E0%B8%B2%E0%B8%A3%E0%B8%AD%E0%B8%A3%E0%B9%88%E0%B8%AD%E0%B8%A2%E0%B9%80%E0%B8%A1%E0%B8%B7%E0%B8%AD%E0%B8%87%E0%B8%AA%E0%B8%B4%E0%B8%87%E0%B8%AB%E0%B9%8C_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201112191540281739'),
    ('ร้านอาหารถ้ำกระแซ',
     'กาญจนบุรี',
     'อาหารไทย',
     'ตำบลลุ่มสุ่ม อำเภอไทรโยค จังหวัดกาญจนบุรี',
     'ร้านอาหารบริเวณถ้ำกระแซ ตำบลลุ่มสุ่ม อำเภอไทรโยค จังหวัดกาญจนบุรี

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=55403163)
ภาพ: ร้านอาหารถ้ำกระแซ Tham Krasae Restaurant - panoramio.jpg — khanawoot / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=55403163 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/9/9b/%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%AD%E0%B8%B2%E0%B8%AB%E0%B8%B2%E0%B8%A3%E0%B8%96%E0%B9%89%E0%B8%B3%E0%B8%81%E0%B8%A3%E0%B8%B0%E0%B9%81%E0%B8%8B_Tham_Krasae_Restaurant_-_panoramio.jpg/960px-%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%AD%E0%B8%B2%E0%B8%AB%E0%B8%B2%E0%B8%A3%E0%B8%96%E0%B9%89%E0%B8%B3%E0%B8%81%E0%B8%A3%E0%B8%B0%E0%B9%81%E0%B8%8B_Tham_Krasae_Restaurant_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=55403163'),
    ('แม่เฉลียว',
     'ชุมพร',
     'อาหารใต้',
     '35 2 ตำบลขันเงิน อำเภอหลังสวน จังหวัดชุมพร 86110',
     'ร้านอาหารไทยและอาหารพื้นเมืองภาคใต้ในอำเภอหลังสวน จังหวัดชุมพร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20200910155310471766)
ภาพ: MaeChaleaw restaurant.jpg — Supisara.plus / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=68640850 (แสดงเป็นภาพย่อ)',
     9.9539054,
     99.07549960000001,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/a/a6/MaeChaleaw_restaurant.jpg/960px-MaeChaleaw_restaurant.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20200910155310471766'),
    ('ร้านอาหารปันหยี',
     'พังงา',
     'อาหารทะเล',
     '78/7 ม.2 อำเภอเมืองพังงา จังหวัดพังงา 82000',
     'ร้านอาหารทะเลที่เกาะปันหยี อำเภอเมืองพังงา จังหวัดพังงา

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20111207113520609)
ภาพ: Paynee Restaurante at the floating muslim village. - panoramio.jpg — Jaakko H. / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=52215779 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/8/89/Paynee_Restaurante_at_the_floating_muslim_village._-_panoramio.jpg/960px-Paynee_Restaurante_at_the_floating_muslim_village._-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20111207113520609'),
    ('The Mountain Grand View',
     'สุราษฎร์ธานี',
     'อาหารไทย',
     'อำเภอเกาะสมุย จังหวัดสุราษฎร์ธานี 84140',
     'ร้านอาหารไทยและเครื่องดื่มบนเกาะสมุย จังหวัดสุราษฎร์ธานี

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20180219161649340650)
ภาพ: The mountain grand view restaurant koh samui - panoramio (1).jpg — Dirk Enthoven / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=56864156 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/e/eb/The_mountain_grand_view_restaurant_koh_samui_-_panoramio_%281%29.jpg/960px-The_mountain_grand_view_restaurant_koh_samui_-_panoramio_%281%29.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20180219161649340650'),
    ('ขนมจีบสังขยาป้าพิณ',
     'ตรัง',
     'อาหารไทย',
     '30/2 หมู่ 10 ถนนเพชรเกษม ตำบลลำภูรา อำเภอห้วยยอด จังหวัดตรัง',
     'ร้านขนมจีบไส้สังขยาพื้นถิ่นตรังในตำบลลำภูรา อำเภอห้วยยอด

ที่มาข้อมูล: กรมส่งเสริมอุตสาหกรรม (IPC10) (https://ipc10.dip.go.th/uploadcontent/IPC10/Purchase%20NEW/ipc10_2564.pdf)
ภาพ: Lam Phura, Trang, Thailand.jpg — Sharon Hahn Darlin / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=163838186 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/c/ca/Lam_Phura%2C_Trang%2C_Thailand.jpg/960px-Lam_Phura%2C_Trang%2C_Thailand.jpg',
     'กรมส่งเสริมอุตสาหกรรม (IPC10)',
     'https://ipc10.dip.go.th/uploadcontent/IPC10/Purchase%20NEW/ipc10_2564.pdf'),
    ('ก๋วยเตี๋ยวปากหม้อ พนมสารคาม บางปะอิน',
     'พระนครศรีอยุธยา',
     'อาหารไทย',
     'ทางหลวงเอเชียหมายเลข 1 ตำบลบ้านหว้า อำเภอบางปะอิน จังหวัดพระนครศรีอยุธยา 13160',
     'ร้านก๋วยเตี๋ยวปากหม้อริมถนนสายเอเชียในอำเภอบางปะอิน จังหวัดพระนครศรีอยุธยา

ที่มาข้อมูล: กระทรวงการท่องเที่ยวและกีฬา (Thailand Tourism Directory) (https://www.thailandtourismdirectory.go.th/restaurants/84040)
ภาพ: Bang Pa-in - Kuai tiao pak mo - ก๋วยเตี๋ยวปากหม้อ พนมสารคาม บางปะอิน (2025) - img 01.jpg — Chainwit. / CC BY 4.0 (https://creativecommons.org/licenses/by/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=159755058 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/6/6b/Bang_Pa-in_-_Kuai_tiao_pak_mo_-_%E0%B8%81%E0%B9%8B%E0%B8%A7%E0%B8%A2%E0%B9%80%E0%B8%95%E0%B8%B5%E0%B9%8B%E0%B8%A2%E0%B8%A7%E0%B8%9B%E0%B8%B2%E0%B8%81%E0%B8%AB%E0%B8%A1%E0%B9%89%E0%B8%AD_%E0%B8%9E%E0%B8%99%E0%B8%A1%E0%B8%AA%E0%B8%B2%E0%B8%A3%E0%B8%84%E0%B8%B2%E0%B8%A1_%E0%B8%9A%E0%B8%B2%E0%B8%87%E0%B8%9B%E0%B8%B0%E0%B8%AD%E0%B8%B4%E0%B8%99_%282025%29_-_img_01.jpg/960px-Bang_Pa-in_-_Kuai_tiao_pak_mo_-_%E0%B8%81%E0%B9%8B%E0%B8%A7%E0%B8%A2%E0%B9%80%E0%B8%95%E0%B8%B5%E0%B9%8B%E0%B8%A2%E0%B8%A7%E0%B8%9B%E0%B8%B2%E0%B8%81%E0%B8%AB%E0%B8%A1%E0%B9%89%E0%B8%AD_%E0%B8%9E%E0%B8%99%E0%B8%A1%E0%B8%AA%E0%B8%B2%E0%B8%A3%E0%B8%84%E0%B8%B2%E0%B8%A1_%E0%B8%9A%E0%B8%B2%E0%B8%87%E0%B8%9B%E0%B8%B0%E0%B8%AD%E0%B8%B4%E0%B8%99_%282025%29_-_img_01.jpg',
     'กระทรวงการท่องเที่ยวและกีฬา (Thailand Tourism Directory)',
     'https://www.thailandtourismdirectory.go.th/restaurants/84040'),
    ('ร้านกาแฟชายทุ่ง',
     'ปทุมธานี',
     'คาเฟ่',
     'คลอง 4 ถนนรังสิต-นครนายก ตำบลบึงยี่โถ อำเภอธัญบุรี จังหวัดปทุมธานี',
     'ร้านกาแฟและเบเกอรี่ริมแปลงนาสาธิตในตำบลบึงยี่โถ อำเภอธัญบุรี จังหวัดปทุมธานี

ที่มาข้อมูล: สำนักงานประชาสัมพันธ์จังหวัดปทุมธานี (https://pathumthani.prd.go.th/th/content/category/detail/id/33/iid/339375)
ภาพ: Chai Thung Coffee by Tris T7.jpg — Tris T7 / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=79808335 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/c/cf/Chai_Thung_Coffee_by_Tris_T7.jpg/960px-Chai_Thung_Coffee_by_Tris_T7.jpg',
     'สำนักงานประชาสัมพันธ์จังหวัดปทุมธานี',
     'https://pathumthani.prd.go.th/th/content/category/detail/id/33/iid/339375'),
    ('Café Amazon พลายชุมพล',
     'พิษณุโลก',
     'คาเฟ่',
     'ตำบลพลายชุมพล อำเภอเมืองพิษณุโลก จังหวัดพิษณุโลก',
     'ร้านกาแฟ Café Amazon ในพื้นที่พลายชุมพล อำเภอเมืองพิษณุโลก

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=99645573)
ภาพ: Cafe Amazon, Phlai Chumphon, Phitsanulok.jpg — SietsL / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=99645573 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/a/a6/Cafe_Amazon%2C_Phlai_Chumphon%2C_Phitsanulok.jpg/960px-Cafe_Amazon%2C_Phlai_Chumphon%2C_Phitsanulok.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=99645573'),
    ('Café Amazon สตาร์พลาซ่า อรัญประเทศ',
     'สระแก้ว',
     'คาเฟ่',
     'ศูนย์การค้าสตาร์พลาซ่า บ้านคลองลึก อำเภออรัญประเทศ จังหวัดสระแก้ว',
     'ร้านกาแฟ Café Amazon ที่ศูนย์การค้าสตาร์พลาซ่า ย่านบ้านคลองลึก อำเภออรัญประเทศ

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=82178007)
ภาพ: Café Amazon, Star Plaza, Aranyaprathet, 2018 (01).jpg — Bahnfrend / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=82178007 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/c/cf/Caf%C3%A9_Amazon%2C_Star_Plaza%2C_Aranyaprathet%2C_2018_%2801%29.jpg/960px-Caf%C3%A9_Amazon%2C_Star_Plaza%2C_Aranyaprathet%2C_2018_%2801%29.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=82178007')
)
INSERT INTO restaurants(name,province_id,category_id,address,description,latitude,longitude,image_url,source_name,source_url)
SELECT s.name,p.id,c.id,s.address,s.description,s.latitude::double precision,s.longitude::double precision,s.image_url,s.source_name,s.source_url
FROM additional_places s
JOIN provinces p ON p.name_th=s.province_name
JOIN categories c ON c.name_th=s.category_name
ON CONFLICT(source_name,source_url) WHERE source_name IS NOT NULL AND source_url IS NOT NULL DO NOTHING;

COMMIT;
