import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Hero from '../components/Home/Hero';
import HomestayCard from '../components/Homestay/HomestayCard';
import { Homestay } from '../types';
import { homestayService } from '../services/homestayService';
import { Leaf, Cloud, Camera } from 'lucide-react';

const Home: React.FC = () => {
  const navigate = useNavigate();
  const [homestays, setHomestays] = React.useState<Homestay[]>([]);

  useEffect(() => {
    const fetchHomestays = async () => {
      try {
        const data = await homestayService.getTopHomestays(8);
        setHomestays(data);
      } catch (error) {
        console.error('Error fetching homestays:', error);
      }
    };

    fetchHomestays();
  }, []);

  const handleSearch = (filters: any) => {
    navigate('/homestays', { state: { filters } });
  };

  const handleHomestayClick = (homestayId: number) => {
    navigate(`/homestay/${homestayId}`);
  };

  return (
    <>
      <Hero onSearch={handleSearch} />

      {/* Giới thiệu */}
      <section className="bg-indigo-50 py-16 px-4">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-3xl font-light text-indigo-700 mb-4">Một nơi để dừng lại</h2>
          <p className="text-lg text-gray-700">
            Mỗi hành trình cần một điểm dừng – và tại Mây Lang Thang, chúng tôi tạo ra những điểm dừng bình yên, nơi bạn tìm thấy chính mình giữa mây trời, rừng thông và tiếng chim ban sớm.
          </p>
        </div>
      </section>

      {/* Homestay nổi bật */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <div className="text-center mb-12">
          <h2 className="text-3xl font-bold text-gray-900 mb-4">
            Homestay nổi bật
          </h2>
          <p className="text-lg text-gray-600">
            Khám phá những homestay được yêu thích nhất
          </p>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          {homestays.slice(0, 8).map((homestay) => (
            <HomestayCard
              key={homestay.id}
              homestay={homestay}
              onClick={() => handleHomestayClick(homestay.id)}
            />
          ))}
        </div>
      </div>

      {/* Giá trị của Mây */}
      <section className="bg-white py-20 px-4">
        <div className="max-w-6xl mx-auto">
          <h2 className="text-center text-3xl font-semibold text-indigo-700 mb-12">Vì sao chọn Mây Lang Thang?</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-12 text-center">
            <div className="flex flex-col items-center">
              <Cloud className="text-indigo-600 w-10 h-10 mb-4" />
              <h3 className="text-xl font-medium mb-2">Không gian thiền tịnh</h3>
              <p className="text-gray-600">
                Tránh xa đô thị ồn ào, mỗi homestay là một thế giới nhỏ để bạn an trú.
              </p>
            </div>
            <div className="flex flex-col items-center">
              <Leaf className="text-indigo-600 w-10 h-10 mb-4" />
              <h3 className="text-xl font-medium mb-2">Gần gũi thiên nhiên</h3>
              <p className="text-gray-600">
                Rừng, suối, đồi và sương – thiên nhiên là người bạn đồng hành tuyệt vời nhất.
              </p>
            </div>
            <div className="flex flex-col items-center">
              <Camera className="text-indigo-600 w-10 h-10 mb-4" />
              <h3 className="text-xl font-medium mb-2">Check-in nghệ thuật</h3>
              <p className="text-gray-600">
                Không gian thơ – từ tường, bàn, ánh sáng đến từng góc nhỏ để bạn lưu lại những khoảnh khắc.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA cuối */}
      <section className="bg-indigo-600 text-white py-20 text-center px-4">
        <div className="max-w-3xl mx-auto">
          <h2 className="text-3xl font-light mb-4">Sẵn sàng cho một hành trình nhẹ tênh?</h2>
          <p className="text-lg mb-6 text-indigo-100">
            Cùng Mây khám phá những nơi bạn chưa từng đến – và những cảm xúc bạn tưởng mình đã quên.
          </p>
          <button
            onClick={() => navigate('/homestays')}
            className="bg-white text-indigo-700 font-medium px-6 py-3 rounded-full shadow-md hover:bg-indigo-50 transition"
          >
            Khám phá Homestay
          </button>
        </div>
      </section>
    </>
  );
};

export default Home;
