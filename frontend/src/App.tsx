import { useQuery } from 'react-query';
import Layout from './components/Layout';
import { Link } from 'react-router-dom';
import { useGetPosts } from './getData';

export interface PostI {
  id: string;
  title: string;
  body: string;
  date: string;
  author: string;
}

function App() {
  const { isSuccess, data, error, isLoading } = useGetPosts();

  if (isLoading)
    return (
      <Layout>
        <div className='flex justify-center items-center'>
          <h3 className=' text-main-white text-xl'>Loading...</h3>;
        </div>
      </Layout>
    );

  return (
    <Layout>
      <>
        <section className='flex text-main-white'>
          {isSuccess &&
            data.slice(0, 1).map((post: PostI) => (
              <section key={post.id} className=' w-1/2 pr-6'>
                <article className=' w-full h-full'>
                  <div className='w-full h-3/5 overflow-hidden'>
                    <img
                      src='https://duet-cdn.vox-cdn.com/thumbor/0x0:4000x2667/640x427/filters:focal(2000x1334:2001x1335):format(webp)/cdn.vox-cdn.com/uploads/chorus_asset/file/23659360/609754850.jpg'
                      alt='rocket'
                      className='w-full h-auto'
                    />
                  </div>
                  <Link to={`post/${post.id}`}>
                    <h1 className=' text-4xl font-extrabold mb-3 w-3/4 ml-8 -mt-4 relative hover:text-main-red'>
                      {post.title}
                    </h1>
                  </Link>
                  <p className='text-2xl font-extralight text-stone-200 ml-8 '>
                    {post.body.split('.')[0]}.
                  </p>
                  <div className='flex mt-6 ml-8 text-lg font-light text-stone-400'>
                    <p className='mr-6 text-main-red text-base'>
                      {post.author}
                    </p>
                    <p className='text-base'>{post.date.split('T')[0]}</p>
                  </div>
                </article>
              </section>
            ))}

          <section className='w-1/2 pl-12'>
            {isSuccess &&
              data.slice(1, 4).map((post: PostI) => (
                <article
                  key={post.id}
                  className='flex items-center w-full mb-10'
                >
                  <img
                    src='https://duet-cdn.vox-cdn.com/thumbor/0x0:6592x4399/1200x800/filters:focal(3296x2200:3297x2201):format(webp)/cdn.vox-cdn.com/uploads/chorus_asset/file/23893495/1238141654.jpg'
                    alt='rocket'
                    className='w-1/3 h-auto'
                  />
                  <section className='pl-6'>
                    <Link to={`post/${post.id}`}>
                      <h1 className=' text-3xl font-bold mb-3 mt-1 relative hover:text-main-red'>
                        {post.title}
                      </h1>
                    </Link>
                    <div className='flex mt-2 text-base font-light text-stone-400'>
                      <p className='mr-6 text-main-red'>{post.author}</p>
                      <p className=''>{post.date.split('T')[0]}</p>
                    </div>
                  </section>
                </article>
              ))}
          </section>
        </section>
      </>
    </Layout>
  );
}

export default App;
