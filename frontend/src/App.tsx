import { useQuery } from 'react-query';
import Layout from './components/Layout';

interface Post {
  id: string;
  title: string;
  body: string;
  date: string;
  author: string;
}

function App() {
  const getPosts = async () => {
    const res = await fetch('http://localhost:3500/posts');
    return res.json();
  };

  const { isSuccess, data, error, isLoading } = useQuery('allPosts', getPosts);

  if (isLoading)
    return (
      <Layout>
        <div className='flex justify-center items-center'>
          <h3 className=' text-main-white text-xl'>Loading...</h3>;
        </div>
      </Layout>
    );
  if (isSuccess) {
    return (
      <Layout>
        <section className='flex text-main-white'>
          {data.map((post: Post) => (
            <article key={post.id} className=' w-1/2'>
              <img
                src='https://duet-cdn.vox-cdn.com/thumbor/0x0:6592x4399/1200x800/filters:focal(3296x2200:3297x2201):format(webp)/cdn.vox-cdn.com/uploads/chorus_asset/file/23893495/1238141654.jpg'
                alt='rocket'
                className='w-full h-auto'
              />
              <h1 className=' text-5xl font-extrabold mb-3 -mt-4 relative hover:text-main-red'>
                {post.title}
              </h1>
              <p className='text-base font-extralight text-stone-200 '>
                {post.body.slice(0, 200)}...
              </p>
              <div className='flex mt-2 text-sm font-light text-stone-400'>
                <p className='mr-6 text-main-red'>{post.author}</p>
                <p className=''>{post.date.split('T')[0]}</p>
              </div>
            </article>
          ))}
        </section>
      </Layout>
    );
  }
}

export default App;
