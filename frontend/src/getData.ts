import { useQuery } from 'react-query';
import { PostI } from './App';

export const useGetPosts = () => {
  const getPosts = async () => {
    const res = await fetch('http://localhost:3500/posts');
    return res.json();
  };

  const { isSuccess, data, error, isLoading, status } = useQuery(
    'allPosts',
    getPosts
  );

  if (data && data.length > 0) {
    data.sort((a: PostI, b: PostI) => {
      return Date.parse(b.date) - Date.parse(a.date);
    });
  }

  console.log(data, status);

  return { isSuccess, data, error, isLoading };
};
