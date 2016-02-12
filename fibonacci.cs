using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using QuickGraph;

namespace Optisynth.Common
{
    using Optisynth.AspectParser;

    public class FibonacciHeap<StorageType> : IPriorityQueue<StorageType>
        where StorageType : IComparable
    {
        //This is the 10th Fibonacci Number
        //private static readonly int MAX_INSERTS_NO_CONSOLIDATE = 89;
        private static readonly int[] FIBO_MATRIX = { 2, 3, 5, 8, 13, 21, 
												      34, 55, 89, 144, 233, 
												      377, 610, 987, 1597, 2584, 
												      4181, 6765, 10946, 17711, 28657, 
												      46368, 75025, 121393, 196418, 317811, 
												      514229, 832040, 1346269, 2178309, 3524578, 
												      5702887, 9227465, 14930352, 24157817, 39088169, 
												      63245986, 102334155, 165580141, 267914296, 433494437, 
												      701408733, 1134903170, 1836311903 };
        private int fiboIndex = 0;
        private int fiboTarget = FIBO_MATRIX[0];
        private int oldFiboTarget = FIBO_MATRIX[0];

        private FibonacciHeapNode<StorageType> minimum = null;
        private readonly StorageType minimumKeyValue;
        private int size = 0;

        public FibonacciHeap(StorageType minimumKeyValue)
        {
            this.minimumKeyValue = minimumKeyValue;
        }

        //[AspectMethod(
        public override string ToString()
        {
            if (IsEmpty)
            {
                return "FibonacciHeap=[]";
            }

            StringBuilder returnBuffer = new StringBuilder(1024);
            returnBuffer.Append("FibonacciHeap={\n");
            printNode(minimum, 1, returnBuffer);
            returnBuffer.Append("};");
            return returnBuffer.ToString();
        }

        private void printNode(FibonacciHeapNode<StorageType> currentNode, int tabs, StringBuilder returnBuffer)
        {
            for (int i = 0; i < tabs; ++i) { returnBuffer.Append('\t'); }
            returnBuffer.Append(currentNode).Append('\n');
            if (currentNode.Child != null)
            {
                printNode(currentNode.Child, tabs + 1, returnBuffer);
            }

            for (FibonacciHeapNode<StorageType> nextNode = currentNode.Right; nextNode != currentNode; nextNode = nextNode.Right)
            {
                for (int i = 0; i < tabs; ++i) { returnBuffer.Append('\t'); }
                returnBuffer.Append(nextNode).Append('\n');
                if (nextNode.Child != null)
                {
                    printNode(nextNode.Child, tabs + 1, returnBuffer);
                }
            }
        }

        public StorageType HeapMinimum { get { return (minimum != null ? minimum.Key : minimumKeyValue); } }
        public StorageType PopMinimum
        {
            get
            {
                StorageType minType = default(StorageType);
                if (minimum != null)
                {
                    minType = minimum.Key;
                    removeMin();
                }
                return minType;
            }
        }
        public IQueueNode<StorageType> HeapNodeMinimum { get { return minimum; } }
        public void clear()
        {
            minimum = null;
            size = 0;
        }

        public bool IsEmpty { get { return minimum == null; } } 
        public int Count { get { return size; } }

        public void decreaseKey(FibonacciHeapNode<StorageType> node, StorageType k)
        {
            if (k.CompareTo(node.Key) > 0)
            {
                throw new Exception("Decrease Key got larger key value");
            }

            node.Key = k;
            FibonacciHeapNode<StorageType> parent = node.Parent;
            if (parent != null && node.Key.CompareTo(parent.Key) < 0)
            {
                cut(node, parent);
                cascadingCut(parent);
            }

            if (node.Key.CompareTo(minimum.Key) < 0)
            {
                minimum = node;
            }
        }

        public void delete(FibonacciHeapNode<StorageType> node)
        {
            decreaseKey(node, minimumKeyValue);
            removeMin();
        }

        public void insert(StorageType key)
        {
            insert(new FibonacciHeapNode<StorageType>(), key);
        }

        public void insert(FibonacciHeapNode<StorageType> node, StorageType key)
        {
            node.Key = key;
            if (minimum != null)
            {
                node.Left = minimum;
                node.Right = minimum.Right;
                minimum.Right = node;
                node.Right.Left = node;

                if (key.CompareTo(minimum.Key) < 0)
                {
                    minimum = node;
                }
            }
            else
            {
                minimum = node;
                minimum.Right = node;
                minimum.Left = node;
            }

            if (++size % fiboTarget == 0)
            {
                oldFiboTarget = fiboIndex;
                fiboTarget = FIBO_MATRIX[++fiboIndex];
                consolidate();
            }
        }

        #region Deprecated removeMin
        ///
        /// The Following method is left commented out, it is the direct although marginally optimized
        ///     Implementation of removeMin() from this psuedo code: http://www.cse.yorku.ca/~aaw/Jason/FibonacciHeapAlgorithm.html
        ///
        /*
        public FibonacciHeapNode<StorageType> removeMin()
        {
            FibonacciHeapNode<StorageType> oldMin = minimum;
            if (oldMin != null)
            {
                int degree = oldMin.Degree;
                FibonacciHeapNode<StorageType> child = oldMin.Child;
                FibonacciHeapNode<StorageType> right;
                while (degree > 0)
                {
                    right = child.Right;

                    child.Left.Right = child.Right;
                    child.Right.Left = child.Left;

                    child.Left = minimum;
                    child.Right = minimum.Right;
                    minimum.Right = child;
                    child.Right.Left = child;

                    child.Parent = null;
                    child = right;
                    --degree;
                }

                minimum.Left.Right = minimum.Right;
                minimum.Right.Left = minimum.Left;

                if (minimum == minimum.Right)
                {
                    minimum = null;
                }
                else
                {
                    minimum = minimum.Right;
                    consolidate();
                }

                --size;
            }
            return oldMin;
        }*/
        #endregion
        /// <summary>
        /// 
        /// </summary>
        /// <returns></returns>
        public FibonacciHeapNode<StorageType> removeMin()
        {
            FibonacciHeapNode<StorageType> oldMin = minimum;
            if (oldMin != null)
            {
                if (oldMin.Degree > 0)
                {
                    FibonacciHeapNode<StorageType> child = oldMin.Child;
                    child.Parent = null;
                    for (FibonacciHeapNode<StorageType> right = child.Right; right != child; right = right.Right)
                    {
                        right.Parent = null;
                    }

                    if (minimum.Right == minimum)
                    {
                        minimum = child;
                        consolidate();
                    }
                    else
                    {
                        FibonacciHeapNode<StorageType> minLeft = minimum.Left;
                        FibonacciHeapNode<StorageType> minRight = minimum.Right;
                        FibonacciHeapNode<StorageType> childLeft = child.Left;
                        FibonacciHeapNode<StorageType> childRight = child.Left.Left;

                        minLeft.Right = childLeft;
                        minRight.Left = childRight;
                        childLeft.Left = minLeft;
                        childRight.Right = minRight;
                        minimum = minRight;

                        consolidate();
                    }
                }
                else if (minimum.Right == minimum)
                {
                    minimum = null;
                }
                else
                {
                    minimum.Left.Right = minimum.Right;
                    minimum.Right.Left = minimum.Left;
                    minimum = minimum.Right;
                    consolidate();
                }

                if (--size / oldFiboTarget == 1 && size % oldFiboTarget == 0)
                {
                    --fiboIndex;
                    fiboTarget = oldFiboTarget;
                    if (fiboIndex > 0)
                    {
                        oldFiboTarget = FIBO_MATRIX[fiboIndex];
                    }
                }
            }

            return oldMin;
        }

        private void cascadingCut(FibonacciHeapNode<StorageType> node)
        {
            FibonacciHeapNode<StorageType> parent = node.Parent;

            if (parent != null)
            {
                if (parent.Marked)
                {
                    cut(node, parent);
                    cascadingCut(parent);
                }
                else
                {
                    parent.Marked = true;
                }
            }
        }
        private void cut(FibonacciHeapNode<StorageType> child, FibonacciHeapNode<StorageType> parent)
        {
            child.Left.Right = child.Right;
            child.Right.Left = child.Left;
            --parent.Degree;

            if (parent.Child == child) { parent.Child = child.Right; }
            if (parent.Degree == 0) { parent.Child = null; }

            child.Parent = null;
            child.Left = minimum;
            child.Right = minimum.Right;
            minimum.Right = child;
            child.Right.Left = child;
            child.Marked = false;
        }
        private void link(FibonacciHeapNode<StorageType> child, FibonacciHeapNode<StorageType> parent)
        {
            ///NOTE: This exception check has been removed it should NEVER occur
            //if (child == parent)
            //{
            //    throw new Exception("Unable to link a node to itself!" + child);
            //}

            child.Left.Right = child.Right;
            child.Right.Left = child.Left;

            child.Parent = parent;
            if (parent.Child != null)
            {
                child.Left = parent.Child;
                child.Right = parent.Child.Right;
                parent.Child.Right = child;
                child.Right.Left = child;
            }
            else
            {
                parent.Child = child;
                child.Left = child;
                child.Right = child;
            }

            ++parent.Degree;
            child.Marked = false;
        }
        private void consolidate()
        {
            Dictionary<int, FibonacciHeapNode<StorageType>> nodeDegreeList = new Dictionary<int, FibonacciHeapNode<StorageType>>();
            
            int numberRoots = 0;
            if (minimum != null)
            {
                numberRoots = 1;
                for (FibonacciHeapNode<StorageType> current = minimum.Right; current != minimum; current = current.Right)
                {
                    ++numberRoots;
                }
            }

            FibonacciHeapNode<StorageType> currentNode = minimum;
            while (numberRoots > 0)
            {
                int degree = currentNode.Degree;
                while (nodeDegreeList.ContainsKey(degree))
                {
                    FibonacciHeapNode<StorageType> child;
                    nodeDegreeList.TryGetValue(degree, out child);
                    if (child == currentNode) {
                        break;
                    }
                    
                    FibonacciHeapNode<StorageType> parent = currentNode;
                    if (currentNode.Key.CompareTo(child.Key) > 0)
                    {
                        FibonacciHeapNode<StorageType> temp = child;
                        child = parent;
                        parent = temp;
                    }

                    if (child == minimum)
                    {
                        minimum = parent;
                    }
                    link(child, parent);
                    currentNode = parent;
                    nodeDegreeList.Remove(degree);
                    ++degree;
                }

                nodeDegreeList[degree] = currentNode;
                currentNode = currentNode.Right;
                --numberRoots;
            }

            FibonacciHeapNode<StorageType> newMin = minimum;
            for (FibonacciHeapNode<StorageType> right = minimum.Right; right != minimum; right = right.Right)
            {
                if (right.Key.CompareTo(newMin.Key) < 0)
                {
                    newMin = right;
                }
            }
            minimum = newMin;
        }

        public IEdgeListGraph<FibonacciHeapNode<StorageType>, FibonacciHeapEdge<StorageType>> ToGraph()
        {
            return new FibonacciGraph<StorageType>(minimum);
        }
    }

    public class FibonacciHeapNode<StorageType> : IQueueNode<StorageType>
        where StorageType : IComparable
    {
        private bool marked = false;
        private int degree = 0;
        private StorageType key;
        private FibonacciHeapNode<StorageType> parent = null;
        private FibonacciHeapNode<StorageType> child = null;
        private FibonacciHeapNode<StorageType> left = null;
        private FibonacciHeapNode<StorageType> right = null;

        public bool Marked { get { return marked; } set { marked = value; } }
        public int Degree { get { return degree; } set { degree = value; } }
        public StorageType Key { get { return key; } set { key = value; } }
        public FibonacciHeapNode<StorageType> Parent { get { return parent; } set { parent = value; } }
        public FibonacciHeapNode<StorageType> Child { get { return child; } set { child = value; } }
        public FibonacciHeapNode<StorageType> Left { get { return left; } set { left = value; } }
        public FibonacciHeapNode<StorageType> Right { get { return right; } set { right = value; } }

        public override string ToString()
        {
            StringBuilder returnBuffer = new StringBuilder();
            returnBuffer.Append("Fibo Node[Degree:").Append(degree).Append(" Key:").Append(key).Append(" Marked:").Append(marked).Append("]");
            return returnBuffer.ToString();
        }
    }

    public class FibonacciHeapEdge<StorageType> : IEdge<FibonacciHeapNode<StorageType>>
        where StorageType : IComparable
    {
        private FibonacciHeapNode<StorageType> source;
        private FibonacciHeapNode<StorageType> target;
        public FibonacciHeapEdge(FibonacciHeapNode<StorageType> source, FibonacciHeapNode<StorageType> target) {
            this.source = source;
            this.target = target;
        }
        public FibonacciHeapNode<StorageType> Source { get { return source; } }
        public FibonacciHeapNode<StorageType> Target { get { return target; } }
    }

    public class FibonacciGraph<StorageType> : IEdgeListGraph<FibonacciHeapNode<StorageType>, FibonacciHeapEdge<StorageType>>
        where StorageType : IComparable
    {
        private IList<FibonacciHeapNode<StorageType>> vertexList = new List<FibonacciHeapNode<StorageType>>();
        private IList<FibonacciHeapEdge<StorageType>> edgeList = new List<FibonacciHeapEdge<StorageType>>();

        public FibonacciGraph(FibonacciHeapNode<StorageType> minimuim)
        {
            recursiveAdd(minimuim);
        }

        private void recursiveAdd(FibonacciHeapNode<StorageType> parentNode)
        {
            for (FibonacciHeapNode<StorageType> currentNode = parentNode.Right; currentNode != parentNode; currentNode = currentNode.Right)
            {
                vertexList.Add(currentNode);
                if (currentNode.Parent != null) { edgeList.Add(new FibonacciHeapEdge<StorageType>(currentNode, currentNode.Parent)); }
                if (currentNode.Left != null) { edgeList.Add(new FibonacciHeapEdge<StorageType>(currentNode, currentNode.Left)); }
                if (currentNode.Right != null) { edgeList.Add(new FibonacciHeapEdge<StorageType>(currentNode, currentNode.Right)); }
                if (currentNode.Child != null)
                {
                    edgeList.Add(new FibonacciHeapEdge<StorageType>(currentNode, currentNode.Child));
                    recursiveAdd(currentNode.Child);
                }
            }
        }

        public bool IsDirected { get { return true; } }
        public bool AllowParallelEdges { get { return false; } }
        public bool IsEdgesEmpty { get { return EdgeCount == 0; } }
        public bool IsVerticesEmpty { get { return VertexCount == 0; } }

        public int EdgeCount { get { return edgeList.Count; } }
        public int VertexCount { get { return vertexList.Count; } }

        public void AddVertex(FibonacciHeapNode<StorageType> vertex)
        {
            vertexList.Add(vertex);
        }
        public void AddEdge(FibonacciHeapEdge<StorageType> edge)
        {
            edgeList.Add(edge);
        }

        public IEnumerable<FibonacciHeapEdge<StorageType>> Edges
        {
            get
            {
                return edgeList;
            }
        }

        public IEnumerable<FibonacciHeapNode<StorageType>> Vertices
        {
            get
            {
                return vertexList;
            }
        }

        public bool ContainsEdge(FibonacciHeapEdge<StorageType> edge)
        {
            return true;
        }
        public bool ContainsVertex(FibonacciHeapNode<StorageType> vertex)
        {
            return true;
        }
    }
}
