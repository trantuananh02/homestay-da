import React from 'react';
import { Cloud, Leaf, Camera, Moon, MapPin, Phone, Mail } from 'lucide-react';
import about from '../asset/about.webp'

const About: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <div className="relative bg-indigo-600 text-white py-20">
        <div className="absolute inset-0 bg-black opacity-20"></div>
        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h1 className="text-4xl md:text-6xl font-bold mb-6">
            Về Mây Lang Thang
          </h1>
          <p className="text-xl max-w-3xl mx-auto opacity-90">
            Một hành trình nhẹ tênh giữa mây trời và những homestay thơ mộng – nơi bạn tìm thấy bình yên giữa những vùng đất diệu kỳ.
          </p>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        {/* Mission Section */}
        <div className="text-center mb-16">
          <h2 className="text-3xl font-bold text-gray-900 mb-6">Tầm nhìn & Sứ mệnh</h2>
          <p className="text-lg text-gray-600 max-w-3xl mx-auto">
            Mây Lang Thang không chỉ là nơi dừng chân, mà là hành trình cảm xúc. 
            Chúng tôi kết nối bạn với những homestay có hồn – nơi thiên nhiên, con người và không gian hòa quyện tạo nên những trải nghiệm đáng nhớ.
          </p>
        </div>

        {/* Values Section */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
          <div className="text-center">
            <div className="bg-indigo-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <Cloud className="h-8 w-8 text-indigo-600" />
            </div>
            <h3 className="text-xl font-semibold mb-2">Nhẹ nhàng & Tự do</h3>
            <p className="text-gray-600">
              Tự do khám phá theo cách riêng của bạn, không gò bó, không ràng buộc – như một đám mây phiêu du.
            </p>
          </div>

          <div className="text-center">
            <div className="bg-indigo-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <Leaf className="h-8 w-8 text-indigo-600" />
            </div>
            <h3 className="text-xl font-semibold mb-2">Gần gũi thiên nhiên</h3>
            <p className="text-gray-600">
              Các homestay nằm giữa rừng, bên đồi, cạnh hồ – nơi bạn có thể cảm nhận hơi thở của đất trời.
            </p>
          </div>

          <div className="text-center">
            <div className="bg-indigo-100 w-16 h-16 rounded-full flex items-center justify-center mx-auto mb-4">
              <Moon className="h-8 w-8 text-indigo-600" />
            </div>
            <h3 className="text-xl font-semibold mb-2">Không gian thiền định</h3>
            <p className="text-gray-600">
              Thiết kế tinh tế, yên tĩnh và gợi cảm hứng giúp bạn nghỉ ngơi, viết lách, thiền hoặc tái tạo năng lượng.
            </p>
          </div>
        </div>

        {/* Story Section */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center mb-16">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-6">Chúng tôi là ai?</h2>
            <div className="space-y-4 text-gray-600">
              <p>
                Mây Lang Thang được tạo nên bởi những tâm hồn yêu cái đẹp, đam mê xê dịch và mong muốn lưu giữ vẻ thơ mộng của những miền đất Việt.
              </p>
              <p>
                Từ một căn homestay nhỏ ở Đà Lạt, chúng tôi đã lan tỏa cảm hứng sống chậm, sống chất đến khắp các cao nguyên, bãi biển, thung lũng và thị trấn mù sương.
              </p>
              <p>
                Đằng sau mỗi homestay là một câu chuyện, là chủ nhà đầy đam mê, là chiếc ghế gỗ nhìn ra rừng thông, là buổi sáng sương rơi trên mái lá – tất cả tạo nên Mây.
              </p>
            </div>
          </div>
          <div className="relative">
            <img
              src={about}
              alt="Chúng tôi là ai"
              className="rounded-2xl shadow-lg w-full h-96 object-cover"
            />
          </div>
        </div>

        {/* Contact Section */}
        <div className="bg-indigo-50 rounded-2xl p-8 text-center">
          <h2 className="text-3xl font-bold text-gray-900 mb-6">Liên hệ với Mây</h2>
          <p className="text-gray-600 mb-8 max-w-2xl mx-auto">
            Hãy để lại lời nhắn nếu bạn cần gợi ý nơi trú chân, tư vấn hành trình, hay đơn giản chỉ là một tách trà giữa mây trời.
          </p>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="flex flex-col items-center">
              <div className="bg-indigo-100 w-12 h-12 rounded-full flex items-center justify-center mb-3">
                <MapPin className="h-6 w-6 text-indigo-600" />
              </div>
              <h3 className="font-semibold mb-1">Địa chỉ</h3>
              <p className="text-gray-600 text-sm">
                45 Triệu Việt Vương, Đà Lạt<br />
                Lâm Đồng, Việt Nam
              </p>
            </div>

            <div className="flex flex-col items-center">
              <div className="bg-indigo-100 w-12 h-12 rounded-full flex items-center justify-center mb-3">
                <Phone className="h-6 w-6 text-indigo-600" />
              </div>
              <h3 className="font-semibold mb-1">Hotline</h3>
              <p className="text-gray-600 text-sm">
                0966 123 456<br />
                (8h - 22h mỗi ngày)
              </p>
            </div>

            <div className="flex flex-col items-center">
              <div className="bg-indigo-100 w-12 h-12 rounded-full flex items-center justify-center mb-3">
                <Mail className="h-6 w-6 text-indigo-600" />
              </div>
              <h3 className="font-semibold mb-1">Email</h3>
              <p className="text-gray-600 text-sm">
                hello@maylangthang.vn<br />
                booking@maylangthang.vn
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default About;
