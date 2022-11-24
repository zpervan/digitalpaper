import Layout from '../components/Layout';
import { useQuery } from 'react-query';
import { useParams } from 'react-router-dom';

const Post = () => {
  let { id } = useParams();

  const getPost = async () => {
    const res = await fetch(`http://localhost:3500/posts/${id}`);
    return res.json();
  };

  const { isSuccess, data, error, isLoading } = useQuery('post', getPost);

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
      <main className='text-white'>
        <h1>{data.title}</h1>
        <div className='flex'>
          <p>{data.author}</p>
          <p>{data.date}</p>
        </div>
        <p>{data.body}</p>
      </main>
    </Layout>
  );
};

export default Post;
