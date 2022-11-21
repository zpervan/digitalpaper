import { useQuery } from 'react-query';
import Layout from './components/Layout';

interface Post {
  id: string;
  title: string;
  body: string;
  date: Date;
  author: string;
}

function App() {
  const getPosts = async () => {
    const res = await fetch('http://localhost:3500/posts');
    return res.json();
  };

  const { data, error, isLoading } = useQuery('allPosts', getPosts);

  if (isLoading) return <h3 className=' text-main-white'>Loading...</h3>;

  return (
    <Layout>
      <section className='flex justify-center'>
        {data.map((el: Post) => {
          <article key={el.id}>
            <h1 className=' font-sourceSerifPro'>{el.title}</h1>
            <div>
              <p>{el.author}</p>
              <p>{el.date.getDate()}</p>
            </div>
            <p>{el.body}</p>
          </article>;
        })}
      </section>
    </Layout>
  );
}

export default App;
