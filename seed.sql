-- Idempotent reference data and 59 real restaurants.
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

-- Additional real places: 23 previously empty provinces (2026-11-13).
-- INSERT-only: existing restaurants, reviews and favorites are preserved.
WITH additional_places_23(name,province_name,category_name,address,description,latitude,longitude,image_url,source_name,source_url) AS (VALUES
    ('จันทรโภชนา สาขาเบญจมราชูทิศ',
     'จันทบุรี',
     'อาหารไทย',
     '102/5-8 อำเภอเมืองจันทบุรี จังหวัดจันทบุรี 22000',
     'ร้านอาหารพื้นเมืองจันทบุรี ภาพแสดงหน้าร้านสาขาเบญจมราชูทิศพร้อมป้ายชื่อจันทรโภชนา

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20111207100227545)
ภาพ: Good Restaurant - panoramio.jpg — Adriaan Castermans / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=53959758 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/2/25/Good_Restaurant_-_panoramio.jpg/960px-Good_Restaurant_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20111207100227545'),
    ('Ploy Samet',
     'ระยอง',
     'อาหารไทย',
     'เกาะเสม็ด จังหวัดระยอง',
     'ร้าน Ploy Samet บนเกาะเสม็ด ภาพเก่าแสดงพื้นที่รับประทานอาหารริมทะเลและสะพานไม้ของร้าน

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=46517329)
ภาพ: Ploy samet 2 - panoramio.jpg — Boonchai C / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=46517329 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/4/46/Ploy_samet_2_-_panoramio.jpg/960px-Ploy_samet_2_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=46517329'),
    ('Benz Thai Restaurant (Soneva Kiri)',
     'ตราด',
     'อาหารไทย',
     'Soneva Kiri เกาะกูด จังหวัดตราด',
     'ร้านอาหารไทย Benz ที่ Soneva Kiri บนเกาะกูด ภาพแสดงระเบียงไม้และบรรยากาศริมน้ำของร้าน

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=55170137)
ภาพ: Ko Kut, Soneva Kiri, Benz Thai Restaurant - panoramio (1).jpg — josef knecht / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=55170137 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/f/f4/Ko_Kut%2C_Soneva_Kiri%2C_Benz_Thai_Restaurant_-_panoramio_%281%29.jpg/960px-Ko_Kut%2C_Soneva_Kiri%2C_Benz_Thai_Restaurant_-_panoramio_%281%29.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=55170137'),
    ('ร้านบะหมี่หยกเป็ดย่าง',
     'ฉะเชิงเทรา',
     'อาหารไทย',
     'บางวัว จังหวัดฉะเชิงเทรา',
     'ร้านบะหมี่หยกเป็ดย่างในย่านบางวัว ภาพแสดงหน้าร้านและจุดเตรียมบะหมี่ เกี๊ยว ปู และหมูแดง

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=57268397)
ภาพ: ร้านบะหมี่หยกเป็ดย่าง - panoramio.jpg — กอบแก้ว เสวีวัลลภ / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=57268397 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/4/4c/%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%9A%E0%B8%B0%E0%B8%AB%E0%B8%A1%E0%B8%B5%E0%B9%88%E0%B8%AB%E0%B8%A2%E0%B8%81%E0%B9%80%E0%B8%9B%E0%B9%87%E0%B8%94%E0%B8%A2%E0%B9%88%E0%B8%B2%E0%B8%87_-_panoramio.jpg/960px-%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%9A%E0%B8%B0%E0%B8%AB%E0%B8%A1%E0%B8%B5%E0%B9%88%E0%B8%AB%E0%B8%A2%E0%B8%81%E0%B9%80%E0%B8%9B%E0%B9%87%E0%B8%94%E0%B8%A2%E0%B9%88%E0%B8%B2%E0%B8%87_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=57268397'),
    ('โอวทึ้ง นายกี๋ (ลอดช่อง)',
     'เพชรบุรี',
     'คาเฟ่',
     'อำเภอเมืองเพชรบุรี จังหวัดเพชรบุรี 76000',
     'ร้านของหวานโอวทึ้ง นายกี๋ ในเมืองเพชรบุรี ภาพแสดงไอศกรีมกาแฟเสิร์ฟในถ้วยแก้วที่ร้าน

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20200929113946475689)
ภาพ: ไอศรีมกาแฟ โอวทึ้งนายกี๋ เพชรบุรี - panoramio.jpg — dusitm / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=56743914 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/0/09/%E0%B9%84%E0%B8%AD%E0%B8%A8%E0%B8%A3%E0%B8%B5%E0%B8%A1%E0%B8%81%E0%B8%B2%E0%B9%81%E0%B8%9F_%E0%B9%82%E0%B8%AD%E0%B8%A7%E0%B8%97%E0%B8%B6%E0%B9%89%E0%B8%87%E0%B8%99%E0%B8%B2%E0%B8%A2%E0%B8%81%E0%B8%B5%E0%B9%8B_%E0%B9%80%E0%B8%9E%E0%B8%8A%E0%B8%A3%E0%B8%9A%E0%B8%B8%E0%B8%A3%E0%B8%B5_-_panoramio.jpg/960px-%E0%B9%84%E0%B8%AD%E0%B8%A8%E0%B8%A3%E0%B8%B5%E0%B8%A1%E0%B8%81%E0%B8%B2%E0%B9%81%E0%B8%9F_%E0%B9%82%E0%B8%AD%E0%B8%A7%E0%B8%97%E0%B8%B6%E0%B9%89%E0%B8%87%E0%B8%99%E0%B8%B2%E0%B8%A2%E0%B8%81%E0%B8%B5%E0%B9%8B_%E0%B9%80%E0%B8%9E%E0%B8%8A%E0%B8%A3%E0%B8%9A%E0%B8%B8%E0%B8%A3%E0%B8%B5_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20200929113946475689'),
    ('ร้านชมจันทร์ซีฟู้ด เกาะยอ',
     'สงขลา',
     'อาหารทะเล',
     'เกาะยอ จังหวัดสงขลา',
     'ร้านอาหารทะเลชมจันทร์ซีฟู้ดที่เกาะยอ ภาพแสดงบรรยากาศริมน้ำจากพื้นที่หลังคามุงจากของร้าน

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=55992534)
ภาพ: ร้านชมจันทร์ซีฟู้ด เกาะยอ - panoramio.jpg — กิตติ เลขะกุล / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=55992534 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/7/72/%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%8A%E0%B8%A1%E0%B8%88%E0%B8%B1%E0%B8%99%E0%B8%97%E0%B8%A3%E0%B9%8C%E0%B8%8B%E0%B8%B5%E0%B8%9F%E0%B8%B9%E0%B9%89%E0%B8%94_%E0%B9%80%E0%B8%81%E0%B8%B2%E0%B8%B0%E0%B8%A2%E0%B8%AD_-_panoramio.jpg/960px-%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B8%8A%E0%B8%A1%E0%B8%88%E0%B8%B1%E0%B8%99%E0%B8%97%E0%B8%A3%E0%B9%8C%E0%B8%8B%E0%B8%B5%E0%B8%9F%E0%B8%B9%E0%B9%89%E0%B8%94_%E0%B9%80%E0%B8%81%E0%B8%B2%E0%B8%B0%E0%B8%A2%E0%B8%AD_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=55992534'),
    ('ก๋วยเตี๋ยวนายฮวด บางปลาม้า',
     'สุพรรณบุรี',
     'อาหารไทย',
     'อำเภอบางปลาม้า จังหวัดสุพรรณบุรี',
     'ร้านก๋วยเตี๋ยวนายฮวดในอำเภอบางปลาม้า ภาพแสดงหน้าร้านพร้อมป้ายชื่อภาษาไทยขนาดใหญ่

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=54877907)
ภาพ: ร้าน ก๋วยเตี๋ยวนายฮวด บางปลาม้า สุพรรณบุรี - panoramio.jpg — samunchon khonsamun / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=54877907 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/5/5c/%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99_%E0%B8%81%E0%B9%8B%E0%B8%A7%E0%B8%A2%E0%B9%80%E0%B8%95%E0%B8%B5%E0%B9%8B%E0%B8%A2%E0%B8%A7%E0%B8%99%E0%B8%B2%E0%B8%A2%E0%B8%AE%E0%B8%A7%E0%B8%94_%E0%B8%9A%E0%B8%B2%E0%B8%87%E0%B8%9B%E0%B8%A5%E0%B8%B2%E0%B8%A1%E0%B9%89%E0%B8%B2_%E0%B8%AA%E0%B8%B8%E0%B8%9E%E0%B8%A3%E0%B8%A3%E0%B8%93%E0%B8%9A%E0%B8%B8%E0%B8%A3%E0%B8%B5_-_panoramio.jpg/960px-%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99_%E0%B8%81%E0%B9%8B%E0%B8%A7%E0%B8%A2%E0%B9%80%E0%B8%95%E0%B8%B5%E0%B9%8B%E0%B8%A2%E0%B8%A7%E0%B8%99%E0%B8%B2%E0%B8%A2%E0%B8%AE%E0%B8%A7%E0%B8%94_%E0%B8%9A%E0%B8%B2%E0%B8%87%E0%B8%9B%E0%B8%A5%E0%B8%B2%E0%B8%A1%E0%B9%89%E0%B8%B2_%E0%B8%AA%E0%B8%B8%E0%B8%9E%E0%B8%A3%E0%B8%A3%E0%B8%93%E0%B8%9A%E0%B8%B8%E0%B8%A3%E0%B8%B5_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=54877907'),
    ('บ้านกล้วย&ไข่คาเฟ่',
     'ลพบุรี',
     'คาเฟ่',
     'ตำบลหนองบัว อำเภอพัฒนานิคม จังหวัดลพบุรี 15140',
     'คาเฟ่บ้านกล้วย&ไข่ในจังหวัดลพบุรี ภาพแสดงประติมากรรมพวงกล้วยขนาดใหญ่ในบริเวณสวนของคาเฟ่

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20260305094630001)
ภาพ: Banana and Egg Cafe in Lopburi.jpg — VN.NguyenDucDuy / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=192376137 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://upload.wikimedia.org/wikipedia/commons/f/f1/Banana_and_Egg_Cafe_in_Lopburi.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20260305094630001'),
    ('The Season by The Cliff Lipe',
     'สตูล',
     'อาหารทะเล',
     'The Cliff Lipe Resort เกาะหลีเป๊ะ จังหวัดสตูล',
     'ร้าน The Season ที่ The Cliff Lipe บนเกาะหลีเป๊ะ ภาพแสดงหอยเชลล์พร้อมเครื่องดื่มและวิวทะเลจากโต๊ะอาหารของร้าน

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=163404562)
ภาพ: Koh Lipe, Adang-Rawi Archipelago, Thailand.jpg — Sharon Hahn Darlin / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=163404562 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/d/d5/Koh_Lipe%2C_Adang-Rawi_Archipelago%2C_Thailand.jpg/960px-Koh_Lipe%2C_Adang-Rawi_Archipelago%2C_Thailand.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=163404562'),
    ('Bangkok Tree House Restaurant',
     'สมุทรปราการ',
     'อาหารไทย',
     '60 หมู่ 1 ตำบลบางน้ำผึ้ง อำเภอพระประแดง จังหวัดสมุทรปราการ 10130',
     'ร้านอาหารของ Bangkok Tree House ในบางน้ำผึ้ง ภาพแสดงทางเข้าพร้อมป้ายชื่อ Bangkok Tree House และป้าย Restaurant & Drink

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20190121141129415613)
ภาพ: Bangkok Tree house - panoramio.jpg — mohigan / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=60552669 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/a/aa/Bangkok_Tree_house_-_panoramio.jpg/960px-Bangkok_Tree_house_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20190121141129415613'),
    ('นิวทะเลเผา',
     'ลำปาง',
     'อาหารทะเล',
     '257/1 ถนนทิพย์ช้าง ตำบลสวนดอก อำเภอเมืองลำปาง จังหวัดลำปาง 52100',
     'ร้านอาหารทะเลบนถนนทิพย์ช้างในตัวเมืองลำปาง มีเมนูอาหารไทยและอาหารทะเล

ที่มาข้อมูล: สำนักงานส่งเสริมการจัดประชุมและนิทรรศการ (TCEB) – Thai MICE Connect (https://www.thaimiceconnect.com/business/detail/16743)
ภาพ: Seafood Restaurant at Night - Lampang - Thailand (34406505503).jpg — Adam Jones from Kelowna, BC, Canada / CC BY-SA 2.0 (https://creativecommons.org/licenses/by-sa/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=64174734 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/3/39/Seafood_Restaurant_at_Night_-_Lampang_-_Thailand_%2834406505503%29.jpg/960px-Seafood_Restaurant_at_Night_-_Lampang_-_Thailand_%2834406505503%29.jpg',
     'สำนักงานส่งเสริมการจัดประชุมและนิทรรศการ (TCEB) – Thai MICE Connect',
     'https://www.thaimiceconnect.com/business/detail/16743'),
    ('The Terracotta Garden at Lamphun (สวนไม้ไทยพ่อเลี้ยงหมื่น)',
     'ลำพูน',
     'คาเฟ่',
     '95 หมู่ 3 ตำบลริมปิง อำเภอเมืองลำพูน จังหวัดลำพูน 51000',
     'คาเฟ่และร้านอาหารในสวนไม้ไทยพ่อเลี้ยงหมื่น ตำบลริมปิง จัดแสดงงานดินเผาและสถาปัตยกรรมที่ได้แรงบันดาลใจจากศิลปะโบราณ

ที่มาข้อมูล: กระทรวงการท่องเที่ยวและกีฬา – Thailand Tourism Directory (https://www.thailandtourismdirectory.go.th/en/attraction/5398)
ภาพ: Terracotta Garden Lamphun in te restaurant garden IMG20251203123159.jpg — FredTC / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=179351453 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/0/00/Terracotta_Garden_Lamphun_in_te_restaurant_garden_IMG20251203123159.jpg/960px-Terracotta_Garden_Lamphun_in_te_restaurant_garden_IMG20251203123159.jpg',
     'กระทรวงการท่องเที่ยวและกีฬา – Thailand Tourism Directory',
     'https://www.thailandtourismdirectory.go.th/en/attraction/5398'),
    ('ใบเฟิร์น (Fern Restaurant)',
     'แม่ฮ่องสอน',
     'อาหารเหนือ',
     '87 ถนนขุนลุมประพาส ตำบลจองคำ อำเภอเมืองแม่ฮ่องสอน จังหวัดแม่ฮ่องสอน 58000',
     'ร้านอาหารบนถนนขุนลุมประพาสในตัวเมืองแม่ฮ่องสอน มีอาหารไทยและอาหารเหนือ

ที่มาข้อมูล: กระทรวงการท่องเที่ยวและกีฬา – Thailand Tourism Directory (https://www.thailandtourismdirectory.go.th/th/restaurants/80673)
ภาพ: ร้านใบเฟริน - panoramio.jpg — ::::=UT=:::: / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=52879923 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/5/59/%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B9%83%E0%B8%9A%E0%B9%80%E0%B8%9F%E0%B8%A3%E0%B8%B4%E0%B8%99_-_panoramio.jpg/960px-%E0%B8%A3%E0%B9%89%E0%B8%B2%E0%B8%99%E0%B9%83%E0%B8%9A%E0%B9%80%E0%B8%9F%E0%B8%A3%E0%B8%B4%E0%B8%99_-_panoramio.jpg',
     'กระทรวงการท่องเที่ยวและกีฬา – Thailand Tourism Directory',
     'https://www.thailandtourismdirectory.go.th/th/restaurants/80673'),
    ('ครัวออมเงิน',
     'พะเยา',
     'อาหารไทย',
     'จังหวัดพะเยา',
     'ร้านครัวออมเงินในจังหวัดพะเยา ภาพสถานที่แสดงป้ายชื่อร้านบริเวณทางเข้าและพื้นที่รับประทานอาหาร

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=54769056)
ภาพ: ครัวออมเงิน - panoramio.jpg — Thanapat Sakulna / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=54769056 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/5/56/%E0%B8%84%E0%B8%A3%E0%B8%B1%E0%B8%A7%E0%B8%AD%E0%B8%AD%E0%B8%A1%E0%B9%80%E0%B8%87%E0%B8%B4%E0%B8%99_-_panoramio.jpg/960px-%E0%B8%84%E0%B8%A3%E0%B8%B1%E0%B8%A7%E0%B8%AD%E0%B8%AD%E0%B8%A1%E0%B9%80%E0%B8%87%E0%B8%B4%E0%B8%99_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=54769056'),
    ('Dairy Queen เซ็นทรัลนครสวรรค์',
     'นครสวรรค์',
     'คาเฟ่',
     'ศูนย์การค้าเซ็นทรัลนครสวรรค์ 699 หมู่ 10 ถนนพหลโยธิน ตำบลปากน้ำโพ อำเภอเมืองนครสวรรค์ จังหวัดนครสวรรค์ 60000',
     'ร้านไอศกรีมและของหวาน Dairy Queen ภายในศูนย์การค้าเซ็นทรัลนครสวรรค์

ที่มาข้อมูล: Central Pattana (mall address) / Wikimedia Commons (branch identity) (https://www.centralpattana.co.th/en/contact-us/branch-contact)
ภาพ: Dairy Queen restaurant at Central Nakhon Sawan.jpg — Reversemos.sapanaht / CC0 (https://creativecommons.org/publicdomain/zero/1.0/)
ที่มาภาพ: https://commons.wikimedia.org/?curid=146526648 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/4/44/Dairy_Queen_restaurant_at_Central_Nakhon_Sawan.jpg/960px-Dairy_Queen_restaurant_at_Central_Nakhon_Sawan.jpg',
     'Central Pattana (mall address) / Wikimedia Commons (branch identity)',
     'https://www.centralpattana.co.th/en/contact-us/branch-contact'),
    ('ภัตตาคารท่าอากาศยานขอนแก่น',
     'ขอนแก่น',
     'อาหารไทย',
     'ชั้น 3 ท่าอากาศยานขอนแก่น อำเภอเมืองขอนแก่น จังหวัดขอนแก่น',
     'ภัตตาคารภายในท่าอากาศยานขอนแก่น ตามภาพและคำอธิบายสถานที่เมื่อปี 2561

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=89965707)
ภาพ: Khon Kaen Restaurant.jpg — Bebiezaza / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=89965707 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/9/99/Khon_Kaen_Restaurant.jpg/960px-Khon_Kaen_Restaurant.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=89965707'),
    ('Rustic Barista',
     'นครราชสีมา',
     'คาเฟ่',
     '519 อำเภอเมืองนครราชสีมา จังหวัดนครราชสีมา 30000',
     'ร้านกาแฟ Rustic Barista ในอำเภอเมืองนครราชสีมา

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20220420084920738508)
ภาพ: 2019 02 Rustic Barista Specialty Coffee Korat 01.jpg — Takeaway / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=76793999 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/6/66/2019_02_Rustic_Barista_Specialty_Coffee_Korat_01.jpg/960px-2019_02_Rustic_Barista_Specialty_Coffee_Korat_01.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20220420084920738508'),
    ('KLIM Restaurant',
     'บุรีรัมย์',
     'อาหารไทย',
     'อำเภอเมืองบุรีรัมย์ จังหวัดบุรีรัมย์',
     'ร้านอาหาร KLIM ในจังหวัดบุรีรัมย์ ตามภาพสถานที่เมื่อปี 2555

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=25300345)
ภาพ: KLIM Restaurant.jpg — Mr.BuriramCN / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=25300345 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/8/85/KLIM_Restaurant.jpg/960px-KLIM_Restaurant.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=25300345'),
    ('goodmook',
     'มุกดาหาร',
     'คาเฟ่',
     '10 ตำบลมุกดาหาร อำเภอเมืองมุกดาหาร จังหวัดมุกดาหาร',
     'คาเฟ่ goodmook ในอำเภอเมืองมุกดาหาร มีเครื่องดื่ม ของหวาน และอาหาร

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20241211142302001)
ภาพ: Good^Mook - panoramio.jpg — Adriaan Castermans / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=53890064 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/c/c7/Good%5EMook_-_panoramio.jpg/960px-Good%5EMook_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=20241211142302001'),
    ('Café Amazon สาขา ปตท.ผานกเค้า',
     'เลย',
     'คาเฟ่',
     'ตำบลผานกเค้า อำเภอภูกระดึง จังหวัดเลย',
     'คาเฟ่ Amazon ในสถานีบริการ ปตท. ผานกเค้า อำเภอภูกระดึง ภาพแสดงระเบียงและสวนของร้านพร้อมฉากหลังผานกเค้า

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=60216299)
ภาพ: ผานกเค้า จากCafe Amazon ใน PTT station - panoramio.jpg — ::::=UT=:::: / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=60216299 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/3/37/%E0%B8%9C%E0%B8%B2%E0%B8%99%E0%B8%81%E0%B9%80%E0%B8%84%E0%B9%89%E0%B8%B2_%E0%B8%88%E0%B8%B2%E0%B8%81Cafe_Amazon_%E0%B9%83%E0%B8%99_PTT_station_-_panoramio.jpg/960px-%E0%B8%9C%E0%B8%B2%E0%B8%99%E0%B8%81%E0%B9%80%E0%B8%84%E0%B9%89%E0%B8%B2_%E0%B8%88%E0%B8%B2%E0%B8%81Cafe_Amazon_%E0%B9%83%E0%B8%99_PTT_station_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=60216299'),
    ('Café Amazon สำโรง อุทุมพรพิสัย',
     'ศรีสะเกษ',
     'คาเฟ่',
     'ตำบลสำโรง อำเภออุทุมพรพิสัย จังหวัดศรีสะเกษ 33120',
     'ร้านกาแฟ Café Amazon ในตำบลสำโรง อำเภออุทุมพรพิสัย จังหวัดศรีสะเกษ ตามภาพสถานที่เมื่อปี 2557

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=58748445)
ภาพ: Samrong, Uthumphon Phisai District, Si Sa Ket 33120, Thailand - panoramio (4).jpg — กิตติ เลขะกุล / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=58748445 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/7/75/Samrong%2C_Uthumphon_Phisai_District%2C_Si_Sa_Ket_33120%2C_Thailand_-_panoramio_%284%29.jpg/960px-Samrong%2C_Uthumphon_Phisai_District%2C_Si_Sa_Ket_33120%2C_Thailand_-_panoramio_%284%29.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=58748445'),
    ('กาแฟกับต้นไม้',
     'นครนายก',
     'คาเฟ่',
     '2 หมู่ 3 อำเภอองครักษ์ จังหวัดนครนายก 26120',
     'ร้านกาแฟกับต้นไม้ในอำเภอองครักษ์ จังหวัดนครนายก มีที่นั่งในบ้านไม้และสวนริมคลอง

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201605181504036944)
ภาพ: มุมกาแฟ Coffee Corner - panoramio.jpg — Thaweesak Churasri / CC BY 3.0 (https://creativecommons.org/licenses/by/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=53409286 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/a/a0/%E0%B8%A1%E0%B8%B8%E0%B8%A1%E0%B8%81%E0%B8%B2%E0%B9%81%E0%B8%9F_Coffee_Corner_-_panoramio.jpg/960px-%E0%B8%A1%E0%B8%B8%E0%B8%A1%E0%B8%81%E0%B8%B2%E0%B9%81%E0%B8%9F_Coffee_Corner_-_panoramio.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=201605181504036944'),
    ('Tiao Tao Tan',
     'นนทบุรี',
     'อาหารไทย',
     'ถนนประชาชื่น ตำบลบ้านใหม่ อำเภอปากเกร็ด จังหวัดนนทบุรี',
     'ร้านอาหารฮาลาลริมคลองบนถนนประชาชื่น ตำบลบ้านใหม่ อำเภอปากเกร็ด จังหวัดนนทบุรี

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=155292615)
ภาพ: Tiao Tao Tan.jpg — Phoebus 28 / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=155292615 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://upload.wikimedia.org/wikipedia/commons/a/a0/Tiao_Tao_Tan.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=155292615')
)
INSERT INTO restaurants(name,province_id,category_id,address,description,latitude,longitude,image_url,source_name,source_url)
SELECT s.name,p.id,c.id,s.address,s.description,s.latitude::double precision,s.longitude::double precision,s.image_url,s.source_name,s.source_url
FROM additional_places_23 s
JOIN provinces p ON p.name_th=s.province_name
JOIN categories c ON c.name_th=s.category_name
ON CONFLICT(source_name,source_url) WHERE source_name IS NOT NULL AND source_url IS NOT NULL DO NOTHING;

-- Six additional real places in existing provinces (2026-11-13).
-- INSERT-only: existing restaurants, reviews and favorites are preserved.
WITH additional_existing_places(name,province_name,category_name,address,description,latitude,longitude,image_url,source_name,source_url) AS (VALUES
    ('Anchan Vegetarian Restaurant',
     'เชียงใหม่',
     'อาหารมังสวิรัติ',
     'ย่านนิมมานเหมินท์ อำเภอเมืองเชียงใหม่ จังหวัดเชียงใหม่',
     'ร้านอาหารมังสวิรัติในเชียงใหม่ มีเมนูอาหารไทยจากผักและเต้าหู้ ภาพแสดงยำหัวปลีทอดของร้าน

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=30500981)
ภาพ: Yam huapli thot 02.jpg — Takeaway / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=30500981 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/5/57/Yam_huapli_thot_02.jpg/960px-Yam_huapli_thot_02.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=30500981'),
    ('WhatCat Cafe',
     'เชียงราย',
     'คาเฟ่',
     'ถนนธนาลัย อำเภอเมืองเชียงราย จังหวัดเชียงราย',
     'คาเฟ่แมวบนถนนธนาลัยในเมืองเชียงราย ภาพแสดงหน้าร้านพร้อมป้าย WhatCat Cafe

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=108582862)
ภาพ: WhatCat Cafe, Tanalai Rd., Chiang Rai, Thailand.jpg — Tree Lightbulb / CC0 (https://creativecommons.org/publicdomain/zero/1.0/)
ที่มาภาพ: https://commons.wikimedia.org/?curid=108582862 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/b/b8/WhatCat_Cafe%2C_Tanalai_Rd.%2C_Chiang_Rai%2C_Thailand.jpg/960px-WhatCat_Cafe%2C_Tanalai_Rd.%2C_Chiang_Rai%2C_Thailand.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=108582862'),
    ('Awesome Coffee Shop',
     'นครราชสีมา',
     'คาเฟ่',
     'อำเภอเมืองนครราชสีมา จังหวัดนครราชสีมา',
     'ร้านกาแฟ Awesome Coffee ในตัวเมืองนครราชสีมา ภาพแสดงพื้นที่นั่งและเคาน์เตอร์ภายในร้าน

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=76794825)
ภาพ: 2019 02 Awesome Coffee Shop in Korat.jpg — Takeaway / CC BY-SA 4.0 (https://creativecommons.org/licenses/by-sa/4.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=76794825 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/3/3f/2019_02_Awesome_Coffee_Shop_in_Korat.jpg/960px-2019_02_Awesome_Coffee_Shop_in_Korat.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=76794825'),
    ('ร้านเจ๊ไฝ',
     'กรุงเทพมหานคร',
     'อาหารไทย',
     'แขวงสำราญราษฎร์ เขตพระนคร กรุงเทพมหานคร',
     'ร้านอาหารไทยริมทางในย่านสำราญราษฎร์ เป็นที่รู้จักจากไข่เจียวปูและผัดขี้เมา

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=136836882)
ภาพ: Bangkok street food - Jay Fai.jpg — Streets of Food / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=136836882 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/a/ab/Bangkok_street_food_-_Jay_Fai.jpg/960px-Bangkok_street_food_-_Jay_Fai.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=136836882'),
    ('Pan and David Restaurant',
     'ชลบุรี',
     'อาหารไทย',
     'เกาะสีชัง จังหวัดชลบุรี',
     'ร้านอาหารบนเกาะสีชัง มีอาหารไทยและอาหารทะเลให้เลือก

ที่มาข้อมูล: Wikimedia Commons (place metadata) (https://commons.wikimedia.org/?curid=54314571)
ภาพ: Pan and David Restaurant - panoramio.jpg — ekeidar / CC BY-SA 3.0 (https://creativecommons.org/licenses/by-sa/3.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=54314571 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/4/4e/Pan_and_David_Restaurant_-_panoramio.jpg/960px-Pan_and_David_Restaurant_-_panoramio.jpg',
     'Wikimedia Commons (place metadata)',
     'https://commons.wikimedia.org/?curid=54314571'),
    ('โกปี้เตี่ยม (Kopitiam by Wilai)',
     'ภูเก็ต',
     'อาหารไทย',
     '18 ถนนถลาง ตำบลตลาดใหญ่ อำเภอเมืองภูเก็ต จังหวัดภูเก็ต 83000',
     'ร้านอาหารบนถนนถลางในเมืองเก่าภูเก็ต เสิร์ฟอาหารบาบ๋า-ย่าหยา อาหารไทย กาแฟและชา

ที่มาข้อมูล: การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1) (https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=2011120109182256)
ภาพ: Kopitiam by Wilai, Phuket (13271249955).jpg — Thanakrit Gu from Shanghai, China / CC BY 2.0 (https://creativecommons.org/licenses/by/2.0)
ที่มาภาพ: https://commons.wikimedia.org/?curid=97443627 (แสดงเป็นภาพย่อ)',
     NULL,
     NULL,
     'https://thumb.wikimedia.org/wikipedia/commons/thumb/f/fb/Kopitiam_by_Wilai%2C_Phuket_%2813271249955%29.jpg/960px-Kopitiam_by_Wilai%2C_Phuket_%2813271249955%29.jpg',
     'การท่องเที่ยวแห่งประเทศไทย (data.go.th restaurant1)',
     'https://datacatalog.tat.or.th/dataset/6d8193a1-7187-4625-a22d-504ab69b427d/resource/53af01a0-9ef5-459a-88c1-befe4690b1f2/download/restaurant.json#PLACE_ID=2011120109182256')
)
INSERT INTO restaurants(name,province_id,category_id,address,description,latitude,longitude,image_url,source_name,source_url)
SELECT s.name,p.id,c.id,s.address,s.description,s.latitude::double precision,s.longitude::double precision,s.image_url,s.source_name,s.source_url
FROM additional_existing_places s
JOIN provinces p ON p.name_th=s.province_name
JOIN categories c ON c.name_th=s.category_name
ON CONFLICT(source_name,source_url) WHERE source_name IS NOT NULL AND source_url IS NOT NULL DO NOTHING;

COMMIT;
